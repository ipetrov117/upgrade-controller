package controller

import (
	"context"
	"fmt"

	lifecyclev1alpha1 "github.com/suse-edge/upgrade-controller/api/v1alpha1"
	"github.com/suse-edge/upgrade-controller/internal/upgrade"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *UpgradePlanReconciler) reconcileHelmChart(ctx context.Context, upgradePlan *lifecyclev1alpha1.UpgradePlan, chart *lifecyclev1alpha1.HelmChart) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	fmt.Printf("-----------------------%s------------------------------------", chart.Name)
	conditionType := lifecyclev1alpha1.GetChartConditionType(chart.PrettyName)
	fmt.Println("conditionType: ", conditionType)
	logger.Info("DependencyCharts: ", len(chart.DependencyCharts))

	if len(chart.DependencyCharts) != 0 {
		for _, depChart := range chart.DependencyCharts {
			fmt.Println("depChart: ", depChart)
			depState, err := r.upgradeHelmChart(ctx, upgradePlan, &depChart)
			if err != nil {
				return ctrl.Result{}, err
			}

			fmt.Println("depState: ", depState)
			if depState != upgrade.ChartStateSucceeded && depState != upgrade.ChartStateVersionAlreadyInstalled {
				setCondition, requeue := evaluateHelmChartState(depState)
				fmt.Println("[30]setCondition: ", setCondition)
				fmt.Println("[31]requeue: ", requeue)
				setCondition(upgradePlan, conditionType, depState.FormattedMessage(depChart.ReleaseName))

				return ctrl.Result{Requeue: requeue}, nil
			}
		}
	}

	coreState, err := r.upgradeHelmChart(ctx, upgradePlan, chart)
	fmt.Println("coreState: ", coreState)
	fmt.Println("err: ", err)
	if err != nil {
		return ctrl.Result{}, err
	}

	if coreState == upgrade.ChartStateNotInstalled && len(chart.DependencyCharts) != 0 {
		setFailedCondition(upgradePlan, conditionType, fmt.Sprintf("'%s' core chart is missing, but dependency charts are present", chart.ReleaseName))
		return ctrl.Result{Requeue: true}, nil
	}

	if coreState != upgrade.ChartStateSucceeded && coreState != upgrade.ChartStateVersionAlreadyInstalled {
		setCondition, requeue := evaluateHelmChartState(coreState)
		fmt.Println("[53]setCondition: ", setCondition)
		fmt.Println("[54]requeue: ", requeue)
		setCondition(upgradePlan, conditionType, coreState.FormattedMessage(chart.ReleaseName))

		return ctrl.Result{Requeue: requeue}, nil
	}

	fmt.Println("AddonCharts: ", len(chart.AddonCharts))
	if len(chart.AddonCharts) != 0 {
		for _, addonChart := range chart.AddonCharts {
			fmt.Println("addonChart: ", addonChart)
			addonState, err := r.upgradeHelmChart(ctx, upgradePlan, &addonChart)
			if err != nil {
				return ctrl.Result{}, err
			}

			fmt.Println("addonState: ", addonState)
			switch addonState {
			case upgrade.ChartStateFailed:
				r.Recorder.Eventf(upgradePlan, corev1.EventTypeWarning, conditionType,
					"'%s' upgraded successfully, but add-on component '%s' failed to upgrade", chart.ReleaseName, addonChart.ReleaseName)
			case upgrade.ChartStateNotInstalled:
				r.Recorder.Eventf(upgradePlan, corev1.EventTypeNormal, conditionType,
					"'%s' add-on component upgrade skipped as it is missing in the cluster", addonChart.ReleaseName)
			case upgrade.ChartStateSucceeded:
				r.Recorder.Eventf(upgradePlan, corev1.EventTypeNormal, conditionType,
					"'%s' add-on component successfully upgraded", addonChart.ReleaseName)
			case upgrade.ChartStateInProgress:
				// mark that current add-on chart upgrade is in progress
				setInProgressCondition(upgradePlan, conditionType, addonState.FormattedMessage(addonChart.ReleaseName))
				return ctrl.Result{Requeue: true}, nil
			case upgrade.ChartStateUnknown:
				return ctrl.Result{}, nil
			}
		}
	}

	// to avoid confusion, when upgrade has been done, use core component message in the component condition
	setCondition, requeue := evaluateHelmChartState(coreState)
	setCondition(upgradePlan, conditionType, coreState.FormattedMessage(chart.ReleaseName))
	fmt.Println("-----------------------------------------------------------")
	return ctrl.Result{Requeue: requeue}, nil
}
