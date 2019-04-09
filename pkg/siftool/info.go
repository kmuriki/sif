// Copyright (c) 2018-2019, Sylabs Inc. All rights reserved.
// Copyright (c) 2018, Divya Cote <divya.cote@gmail.com> All rights reserved.
// Copyright (c) 2017, SingularityWare, LLC. All rights reserved.
// Copyright (c) 2017, Yannick Cote <yhcote@gmail.com> All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package siftool

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sylabs/sif/pkg/sif"
)

// Info implements 'siftool info' sub-command
func Info() *cobra.Command {
	return &cobra.Command{
		Use:   "info <descriptorid> <containerfile>",
		Short: "Display detailed information of object descriptors",
		Args:  cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("while converting input descriptor id: %s", err)
			}

			fimg, err := sif.LoadContainer(args[1], true)
			if err != nil {
				return err
			}
			defer func() {
				if err := fimg.UnloadContainer(); err != nil {
					fmt.Println("Error unloading container: ", err)
				}
			}()

			fmt.Print(fimg.FmtDescrInfo(uint32(id)))

			return nil
		},
		DisableFlagsInUseLine: true,
	}
}
