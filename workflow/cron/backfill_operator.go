package cron

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func (woc *cronWfOperationCtx) backfill() *v1alpha1.BackfillStatus {
	return woc.cronWf.Status.Backfill
}

func (woc *cronWfOperationCtx) operateBackfill(ctx context.Context) error {
	if woc.backfill() == nil {
		// Called in error, probably return an error here
		return nil
	}

	if woc.backfill().Next.After(woc.backfill().Until.Time) {
		// We are done backfilling
		woc.cronWf.Status.Backfill = nil
		woc.persistUpdate(ctx)
		return nil
	}

	switch {
	case woc.backfill().Strategy.Sequential != nil:
		return woc.operateSequentialBackfill()
	case woc.backfill().Strategy.Batch != nil:
		// Do
		return nil
	default:
		return fmt.Errorf("cron workflow backfill strategy is not valid")
	}
}

func (woc *cronWfOperationCtx) operateSequentialBackfill() error {

	if len(woc.backfill().Strategy.Active) == 1 {
		if woc.backfill().Strategy.Active[0].Workflow
	}
}
