package controller

import (
	"context"
	"fmt"

	lifecyclev1alpha1 "github.com/suse-edge/upgrade-controller/api/v1alpha1"
	"github.com/suse-edge/upgrade-controller/internal/upgrade"
	"github.com/suse-edge/upgrade-controller/pkg/release"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func getChartConditionType(prettyName string) string {
	return fmt.Sprintf("%sUpgraded", prettyName)
}

func (r *UpgradePlanReconciler) reconcileHelmChart(ctx context.Context, upgradePlan *lifecyclev1alpha1.UpgradePlan, chart *release.HelmChart) (ctrl.Result, error) {
	conditionType := getChartConditionType(chart.PrettyName)
	if len(chart.DependencyCharts) != 0 {
		for _, depChart := range chart.DependencyCharts {
			depState, err := r.upgradeHelmChart(ctx, upgradePlan, &depChart)
			if err != nil {
				return ctrl.Result{}, err
			}

			if depState != upgrade.ChartStateSucceeded && depState != upgrade.ChartStateVersionAlreadyInstalled {
				setCondition, requeue := evaluateHelmChartState(depState)
				setCondition(upgradePlan, conditionType, depState.FormattedMessage(depChart.ReleaseName))

				return ctrl.Result{Requeue: requeue}, err
			}
		}
	}

	coreState, err := r.upgradeHelmChart(ctx, upgradePlan, chart)
	if err != nil {
		return ctrl.Result{}, err
	}

	if coreState == upgrade.ChartStateNotInstalled && len(chart.DependencyCharts) != 0 {
		setFailedCondition(upgradePlan, conditionType, fmt.Sprintf("%s core chart is missing, although dependency charts are present", chart.ReleaseName))
		return ctrl.Result{Requeue: true}, nil
	}

	if coreState != upgrade.ChartStateSucceeded && coreState != upgrade.ChartStateVersionAlreadyInstalled {
		setCondition, requeue := evaluateHelmChartState(coreState)
		setCondition(upgradePlan, conditionType, coreState.FormattedMessage(chart.ReleaseName))

		return ctrl.Result{Requeue: requeue}, err
	}

	if len(chart.AddonCharts) != 0 {
		for _, addonChart := range chart.AddonCharts {
			addonState, err := r.upgradeHelmChart(ctx, upgradePlan, &addonChart)
			if err != nil {
				return ctrl.Result{}, err
			}

			switch addonState {
			case upgrade.ChartStateFailed:
				msg := fmt.Sprintf("Main component '%s' upgraded successfully, but add-on component '%s' failed to upgrade", chart.ReleaseName, addonChart.ReleaseName)
				r.recordPlanEvent(upgradePlan, corev1.EventTypeWarning, "ChartTest", msg)

				fallthrough
			case upgrade.ChartStateNotInstalled, upgrade.ChartStateVersionAlreadyInstalled:
				msg := fmt.Sprintf("%s add-on component upgrade skipped as it is missing in the cluster", addonChart.ReleaseName)
				r.recordPlanEvent(upgradePlan, corev1.EventTypeNormal, "ChartTest", msg)
			default:
				msg := fmt.Sprintf("%s add-on component successfully upgraded", addonChart.ReleaseName)
				r.recordPlanEvent(upgradePlan, corev1.EventTypeNormal, "ChartTest", msg)
			}
		}
	}
	setCondition, requeue := evaluateHelmChartState(coreState)
	setCondition(upgradePlan, conditionType, coreState.FormattedMessage(chart.ReleaseName))
	return ctrl.Result{Requeue: requeue}, nil
}
