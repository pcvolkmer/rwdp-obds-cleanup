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
	_ "embed"
	"encoding/json"
	"strings"
	"testing"
)

//go:embed tests/test.json
var testRecord []byte

//go:embed tests/fixedtest.json
var fixedTestRecord []byte

func TestShouldDeserializeJson(t *testing.T) {
	record, err := ParseRecordValue(testRecord)

	if err != nil {
		t.Errorf("Cannot deserialize record")
	}

	if !strings.HasPrefix(record.Payload.XmlDaten, "<?xml") {
		t.Errorf("XmlDaten does not contain expected XML data")
	}
}

func TestShouldReplacePatientenIdInXmlPayload(t *testing.T) {
	var record RecordValue

	err := json.Unmarshal(testRecord, &record)
	if err != nil {
		t.Errorf("Cannot deserialize record")
	}

	record.RemoveLeadingZerosFromPatientIds()

	fixedJson, err := record.ToJson()
	if err != nil {
		t.Errorf("Cannot serialize record")
	}

	if strings.Contains(string(fixedJson), `Patient_ID=\"00001234\"`) {
		t.Errorf("Patient_ID was not replaced")
	}

	if string(fixedJson) != string(fixedTestRecord) {
		t.Errorf("Fixed Json does not match expected record")
	}
}

func TestShouldNotModifyRecordWithoutLeadingZeros(t *testing.T) {
	var record RecordValue

	err := json.Unmarshal(fixedTestRecord, &record)
	if err != nil {
		t.Errorf("Cannot deserialize record")
	}

	record.RemoveLeadingZerosFromPatientIds()

	fixedJson, err := record.ToJson()
	if err != nil {
		t.Errorf("Cannot serialize record")
	}

	if strings.Contains(string(fixedJson), `Patient_ID=\"00001234\"`) {
		t.Errorf("Patient_ID was not replaced")
	}

	if string(fixedJson) != string(fixedTestRecord) {
		t.Errorf("Fixed Json does not match expected record")
	}
}
