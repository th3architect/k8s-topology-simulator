/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"

	"k8s.io/klog/v2"
)

// parseResult parses outputData to evaluation metrics and writes back to a
// result file
func parseResult(file string, outputQueue <-chan outputData) (err error) {
	outputFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() {
		cerr := outputFile.Close()
		if cerr != nil {
			klog.Errorf("close output file %s with an error %v", file, cerr)
		}
		if err == nil {
			err = cerr
		}
	}()

	klog.Infof("Writing output to file %v\n", file)
	writer := csv.NewWriter(outputFile)

	title := []string{"input name", "score", "in-zone-traffic score", "deviation score", "slice score", "max deviation", "mean deviation", "SD of deviation"}
	err = writer.Write(title)
	if err != nil {
		return err
	}

	for rowData, more := <-outputQueue; more; rowData, more = <-outputQueue {
		// use in zone traffic percentage to be in zone traffic score
		inZoneTrafficScore := rowData.result.InZoneTraffic * 100
		// use mean deviation to calcualte deviation score
		deviationMaxScore := 100.0 - rowData.result.MaxDeviation*100
		deviationMeanScore := 100.0 - rowData.result.MeanDeviation*100
		deviationScore := 0.5*deviationMaxScore + 0.5*deviationMeanScore
		// use number of EndpointSlices deviation to calculate sliceScore
		numberOfOriginalSlices := math.Ceil(float64(rowData.endpoints) / endpointsPerSlice)
		sliceScore := (numberOfOriginalSlices / float64(rowData.endpointSlices)) * 100
		// calculate total score based on two scores above
		totalScore := inZoneTrafficScoreWeight*inZoneTrafficScore + deviationScoreWeight*deviationScore + sliceScoreWeight*sliceScore

		data := []string{rowData.name}
		if rowData.result.Invalid {
			data = append(data, []string{"invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid"}...)
		} else {
			data = append(data, strconv.FormatFloat(totalScore, 'f', 4, 64))
			data = append(data, strconv.FormatFloat(inZoneTrafficScore, 'f', 4, 64))
			data = append(data, strconv.FormatFloat(deviationScore, 'f', 4, 64))
			data = append(data, strconv.FormatFloat(sliceScore, 'f', 4, 64))
			data = append(data, strconv.FormatFloat(rowData.result.MaxDeviation*100, 'f', 4, 64)+"%")
			data = append(data, strconv.FormatFloat(rowData.result.MeanDeviation*100, 'f', 4, 64)+"%")
			data = append(data, strconv.FormatFloat(rowData.result.DeviationSD, 'f', 4, 64))
		}

		err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	err = writer.Error()
	return err
}
