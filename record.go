/*
 * This file is part of rwdp-obds-cleanup
 *
 * Copyright (c) 2024  Paul-Christian Volkmer
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

// Kafka RecordValue

type RecordValue struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

func ParseRecordValue(recordString []byte) (RecordValue, error) {
	var record RecordValue
	err := json.Unmarshal(recordString, &record)
	if err != nil {
		return RecordValue{}, err
	}

	return record, nil
}

func (record *RecordValue) RemoveLeadingZerosFromPatientIds() {
	patientIdRegEx := regexp.MustCompile(`Patient_ID="[^"]+"`)

	for _, match := range patientIdRegEx.FindAllString(record.Payload.XmlDaten, -1) {
		zeroRegEx := regexp.MustCompile(`Patient_ID="0+`)
		fixed := zeroRegEx.ReplaceAllLiteralString(match, `Patient_ID="`)

		record.Payload.XmlDaten = strings.ReplaceAll(record.Payload.XmlDaten, match, fixed)
	}
}

func (record *RecordValue) ToJson() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(record)

	if err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}

type Schema struct {
	Type     string  `json:"type"`
	Fields   []Field `json:"fields"`
	Optional bool    `json:"optional"`
}

type Field struct {
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Field    string `json:"field"`
}

type Payload struct {
	Year           int    `json:"YEAR"`
	Versionsnummer int    `json:"VERSIONSNUMMER"`
	Id             int    `json:"ID"`
	XmlDaten       string `json:"XML_DATEN"`
}
