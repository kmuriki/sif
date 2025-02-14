// Copyright (c) 2021, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sif

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestFileImage_GetDescriptors(t *testing.T) {
	ds := []Descriptor{
		{
			Datatype: DataPartition,
			Used:     true,
			ID:       1,
			Groupid:  1 | DescrGroupMask,
			Link:     DescrUnusedLink,
		},
		{
			Datatype: DataSignature,
			Used:     true,
			ID:       2,
			Groupid:  1 | DescrGroupMask,
			Link:     1,
		},
		{
			Datatype: DataSignature,
			Used:     true,
			ID:       3,
			Groupid:  DescrUnusedGroup,
			Link:     1 | DescrGroupMask,
		},
	}

	tests := []struct {
		name    string
		fns     []DescriptorSelectorFunc
		wantErr error
		wantIDs []uint32
	}{
		{
			name: "DataPartition",
			fns: []DescriptorSelectorFunc{
				WithDataType(DataPartition),
			},
			wantIDs: []uint32{1},
		},
		{
			name: "DataSignature",
			fns: []DescriptorSelectorFunc{
				WithDataType(DataSignature),
			},
			wantIDs: []uint32{2, 3},
		},
		{
			name: "DataSignatureGroupID",
			fns: []DescriptorSelectorFunc{
				WithDataType(DataSignature),
				WithGroupID(1),
			},
			wantIDs: []uint32{2},
		},
		{
			name: "NoGroupID",
			fns: []DescriptorSelectorFunc{
				WithNoGroup(),
			},
			wantIDs: []uint32{3},
		},
		{
			name: "GroupID",
			fns: []DescriptorSelectorFunc{
				WithGroupID(1),
			},
			wantIDs: []uint32{1, 2},
		},
		{
			name: "GroupIDInvalidGroupID",
			fns: []DescriptorSelectorFunc{
				WithGroupID(0),
			},
			wantErr: ErrInvalidGroupID,
		},
		{
			name: "LinkedID",
			fns: []DescriptorSelectorFunc{
				WithLinkedID(1),
			},
			wantIDs: []uint32{2},
		},
		{
			name: "LinkedIDInvalidObjectID",
			fns: []DescriptorSelectorFunc{
				WithLinkedID(0),
			},
			wantErr: ErrInvalidObjectID,
		},
		{
			name: "LinkedGroupID",
			fns: []DescriptorSelectorFunc{
				WithLinkedGroupID(1),
			},
			wantIDs: []uint32{3},
		},
		{
			name: "LinkedGroupIDInvalidGroupID",
			fns: []DescriptorSelectorFunc{
				WithLinkedGroupID(0),
			},
			wantErr: ErrInvalidGroupID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fimg := &FileImage{descrArr: ds}

			ds, err := fimg.GetDescriptors(tt.fns...)
			if got, want := err, tt.wantErr; !errors.Is(got, want) {
				t.Fatalf("got error %v, want %v", got, want)
			}
			if got, want := len(ds), len(tt.wantIDs); got != want {
				t.Errorf("got %v IDs, want %v", got, want)
			}
			for i := range ds {
				if got, want := ds[i].GetID(), tt.wantIDs[i]; got != want {
					t.Errorf("got ID %v, want %v", got, want)
				}
			}
		})
	}
}

func TestFileImage_GetDescriptor(t *testing.T) {
	extra := Partition{
		Fstype:   FsSquash,
		Parttype: PartPrimSys,
	}
	copy(extra.Arch[:], HdrArch386)

	b := bytes.Buffer{}
	if err := binary.Write(&b, binary.LittleEndian, extra); err != nil {
		t.Fatal(err)
	}

	primPartDescr := Descriptor{
		Datatype: DataPartition,
		Used:     true,
		ID:       1,
		Groupid:  1 | DescrGroupMask,
		Link:     DescrUnusedLink,
	}
	copy(primPartDescr.Extra[:], b.Bytes())

	ds := []Descriptor{
		primPartDescr,
		{
			Datatype: DataSignature,
			Used:     true,
			ID:       2,
			Groupid:  1 | DescrGroupMask,
			Link:     1,
		},
		{
			Datatype: DataSignature,
			Used:     true,
			ID:       3,
			Groupid:  DescrUnusedGroup,
			Link:     1 | DescrGroupMask,
		},
	}

	tests := []struct {
		name    string
		fns     []DescriptorSelectorFunc
		wantErr error
		wantID  uint32
	}{
		{
			name: "ID",
			fns: []DescriptorSelectorFunc{
				WithID(1),
			},
			wantID: 1,
		},
		{
			name: "InvalidObjectID",
			fns: []DescriptorSelectorFunc{
				WithID(0),
			},
			wantErr: ErrInvalidObjectID,
		},
		{
			name: "MultipleObjectsFound",
			fns: []DescriptorSelectorFunc{
				WithGroupID(1),
			},
			wantErr: ErrMultipleObjectsFound,
		},
		{
			name: "ObjectNotFound",
			fns: []DescriptorSelectorFunc{
				WithGroupID(2),
			},
			wantErr: ErrObjectNotFound,
		},
		{
			name: "PartitionType",
			fns: []DescriptorSelectorFunc{
				WithPartitionType(PartPrimSys),
			},
			wantID: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fimg := &FileImage{descrArr: ds}

			d, err := fimg.GetDescriptor(tt.fns...)
			if got, want := err, tt.wantErr; !errors.Is(got, want) {
				t.Fatalf("got error %v, want %v", got, want)
			}
			if got, want := d.GetID(), tt.wantID; got != want {
				t.Errorf("got ID %v, want %v", got, want)
			}
		})
	}
}

func TestFileImage_WithDescriptors(t *testing.T) {
	ds := []Descriptor{
		{
			Datatype: DataPartition,
			Used:     true,
			ID:       1,
			Groupid:  1 | DescrGroupMask,
			Link:     DescrUnusedLink,
		},
		{
			Datatype: DataSignature,
			Used:     true,
			ID:       2,
			Groupid:  DescrUnusedGroup,
			Link:     1 | DescrGroupMask,
		},
		{
			Datatype: DataSignature,
			Used:     false,
			ID:       3,
			Groupid:  DescrUnusedGroup,
			Link:     DescrUnusedLink,
		},
	}

	tests := []struct {
		name string
		fn   func(t *testing.T) func(d Descriptor) bool
	}{
		{
			name: "ReturnTrue",
			fn: func(t *testing.T) func(d Descriptor) bool {
				return func(d Descriptor) bool {
					if id := d.GetID(); id > 1 {
						t.Errorf("unexpected ID: %v", id)
					}
					return true
				}
			},
		},
		{
			name: "ReturnFalse",
			fn: func(t *testing.T) func(d Descriptor) bool {
				return func(d Descriptor) bool {
					if id := d.GetID(); id > 2 {
						t.Errorf("unexpected ID: %v", id)
					}
					return false
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileImage{descrArr: ds}

			f.WithDescriptors(tt.fn(t))
		})
	}
}
