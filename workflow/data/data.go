package data

import (
	"encoding/json"
	"fmt"
	"github.com/antonmedv/expr"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
)

func ProcessData(data *wfv1.Data, processor wfv1.DataSourceProcessor) (interface{}, error) {
	sourcedData, err := processSource(data.Source, processor)
	if err != nil {
		return nil, fmt.Errorf("unable to process data source: %w", err)
	}
	transformedData, err := processTransformation(sourcedData, &data.Transformation)
	if err != nil {
		return nil, fmt.Errorf("unable to process data transformation: %w", err)
	}
	return transformedData, nil
}

func processSource(source *wfv1.DataSource, processor wfv1.DataSourceProcessor) (interface{}, error) {
	if source == nil {
		return nil, fmt.Errorf("no source is used for data template")
	}

	var data interface{}
	var err error
	switch {
	case source.ArtifactPaths != nil:
		data, err = processor.ProcessArtifactPaths(source.ArtifactPaths)
		if err != nil {
			return nil, fmt.Errorf("unable to source artifact paths: %w", err)
		}
	case source.Raw != "":
		err = json.Unmarshal([]byte(source.Raw), &data)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal raw source: %w", err)
		}
	default:
		return nil, fmt.Errorf("no valid source is used for data template")
	}

	return data, nil
}

func processTransformation(data interface{}, transformation *wfv1.Transformation) (interface{}, error) {
	if transformation == nil {
		return data, nil
	}

	var err error
	for i, step := range *transformation {
		preData := data
		switch {
		case step.Expression != "":
			data, err = processExpression(step.Expression, data)
		}
		if err != nil {
			logrus.Debugf("data state at time of error: %+v, type: %T", preData, preData)
			return nil, fmt.Errorf("error processing data step %d: %w", i, err)
		}
	}

	return data, nil
}

func processExpression(expression string, data interface{}) (interface{}, error) {
	return expr.Eval(expression, map[string]interface{}{"data": data})
}
