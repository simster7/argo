package raw

import (
	"os"

	"github.com/simster7/argo/v2/errors"
	wfv1 "github.com/simster7/argo/v2/pkg/apis/workflow/v1alpha1"
)

type RawArtifactDriver struct {
}

// Store raw content as artifact
func (a *RawArtifactDriver) Load(artifact *wfv1.Artifact, path string) error {
	lf, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = lf.Close()
	}()

	_, err = lf.WriteString(artifact.Raw.Data)
	return err
}

// Save is unsupported for raw output artifacts
func (g *RawArtifactDriver) Save(string, *wfv1.Artifact) error {
	return errors.Errorf(errors.CodeBadRequest, "Raw output artifacts unsupported")
}
