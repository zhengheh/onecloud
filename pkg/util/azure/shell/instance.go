package shell

import (
	"context"
	"fmt"

	"yunion.io/x/onecloud/pkg/util/azure"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type InstanceListOptions struct {
		Classic   bool `help:"List classic instance"`
		ScaleSets bool `help:"List Scale Sets instance"`
		Limit     int  `help:"page size"`
		Offset    int  `help:"page offset"`
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *azure.SRegion, args *InstanceListOptions) error {
		if args.Classic {
			instances, err := cli.GetClassicInstances()
			if err != nil {
				return err
			}
			printList(instances, len(instances), args.Offset, args.Limit, []string{})
			return nil
		} else if args.ScaleSets {
			instances, err := cli.GetInstanceScaleSets()
			if err != nil {
				return err
			}
			printList(instances, len(instances), args.Offset, args.Limit, []string{})
			return nil
		}
		instances, err := cli.GetInstances()
		if err != nil {
			return err
		}
		printList(instances, len(instances), args.Offset, args.Limit, []string{})
		return nil
	})

	type InstanceSizeListOptions struct {
		Location string
	}

	shellutils.R(&InstanceSizeListOptions{}, "instance-size-list", "List intances", func(cli *azure.SRegion, args *InstanceSizeListOptions) error {
		if vmSize, err := cli.GetVMSize(args.Location); err != nil {
			return err
		} else {
			printObject(vmSize)
			return nil
		}
	})
	shellutils.R(&InstanceSizeListOptions{}, "resource-sku-list", "List resource sku", func(cli *azure.SRegion, args *InstanceSizeListOptions) error {
		skus, err := cli.GetResourceSkus(args.Location)
		if err != nil {
			return err
		}
		printList(skus, len(skus), 0, 0, []string{})
		return nil
	})

	type InstanceCreateOptions struct {
		NAME          string `help:"Name of instance"`
		IMAGE         string `help:"image ID"`
		CPU           int8   `help:"CPU count"`
		MEMORYGB      int    `help:"MemoryGB"`
		SYSDISKSIZEGB int    `help:"System Disk Size"`
		Disk          []int  `help:"Data disk sizes int GB"`
		STORAGE       string `help:"Storage type"`
		NETWORK       string `help:"Network ID"`
		PASSWD        string `help:"password"`
		PublicKey     string `help:"PublicKey"`
		OsType        string `help:"Operation system type" choices:"Linux|Windows"`
	}
	shellutils.R(&InstanceCreateOptions{}, "instance-create", "Create a instance", func(cli *azure.SRegion, args *InstanceCreateOptions) error {
		instance, e := cli.CreateInstanceSimple(args.NAME, args.IMAGE, args.OsType, args.CPU, args.MEMORYGB, args.SYSDISKSIZEGB, args.STORAGE, args.Disk, args.NETWORK, args.PASSWD, args.PublicKey)
		if e != nil {
			return e
		}
		printObject(instance)
		return nil
	})

	type InstanceOptions struct {
		ID string `help:"Instance ID"`
	}
	shellutils.R(&InstanceOptions{}, "instance-show", "Show intance detail", func(cli *azure.SRegion, args *InstanceOptions) error {
		if instance, err := cli.GetInstance(args.ID); err != nil {
			return err
		} else {
			printObject(instance)
			return nil
		}
	})

	shellutils.R(&InstanceOptions{}, "instance-start", "Start intance", func(cli *azure.SRegion, args *InstanceOptions) error {
		return cli.StartVM(args.ID)
	})

	shellutils.R(&InstanceOptions{}, "instance-delete", "Delete intance", func(cli *azure.SRegion, args *InstanceOptions) error {
		return cli.DeleteVM(args.ID)
	})

	shellutils.R(&InstanceOptions{}, "instance-stop", "Stop intance", func(cli *azure.SRegion, args *InstanceOptions) error {
		return cli.StopVM(args.ID, true)
	})

	type InstanceRebuildOptions struct {
		ID        string `help:"Instance ID"`
		CPU       int8   `help:"Instance CPU core"`
		MEMORYMB  int    `help:"Instance Memory MB"`
		IMAGE     string `help:"Image ID"`
		Password  string `help:"pasword"`
		PublicKey string `help:"Public Key"`
		Size      int    `help:"system disk size in GB"`
	}
	shellutils.R(&InstanceRebuildOptions{}, "instance-rebuild-root", "Reinstall virtual server system image", func(cli *azure.SRegion, args *InstanceRebuildOptions) error {
		instance, err := cli.GetInstance(args.ID)
		if err != nil {
			return err
		}
		diskId, err := cli.ReplaceSystemDisk(instance, args.CPU, args.MEMORYMB, args.IMAGE, args.Password, args.PublicKey, args.Size)
		if err != nil {
			return err
		}
		fmt.Printf("New diskId is %s", diskId)
		return nil
	})

	type InstanceDiskOptions struct {
		ID   string `help:"Instance ID"`
		DISK string `help:"Disk ID"`
	}
	shellutils.R(&InstanceDiskOptions{}, "instance-attach-disk", "Attach a disk to intance", func(cli *azure.SRegion, args *InstanceDiskOptions) error {
		return cli.AttachDisk(args.ID, args.DISK)
	})

	shellutils.R(&InstanceDiskOptions{}, "instance-detach-disk", "Attach a disk to intance", func(cli *azure.SRegion, args *InstanceDiskOptions) error {
		return cli.DetachDisk(args.ID, args.DISK)
	})

	type InstanceConfigOptions struct {
		ID     string `help:"Instance ID"`
		NCPU   int    `help:"Number of cpu core"`
		MEMERY int    `helo:"Instance memery in mb"`
	}

	shellutils.R(&InstanceConfigOptions{}, "instance-change-conf", "Attach a disk to intance", func(cli *azure.SRegion, args *InstanceConfigOptions) error {
		return cli.ChangeVMConfig(context.Background(), args.ID, args.NCPU, args.MEMERY)
	})

	type InstanceDeployOptions struct {
		ID        string `help:"Instance ID"`
		Password  string `help:"Password for instance"`
		PublicKey string `helo:"Deploy ssh_key for instance"`
	}

	shellutils.R(&InstanceDeployOptions{}, "instance-reset-password", "Reset intance password", func(cli *azure.SRegion, args *InstanceDeployOptions) error {
		return cli.DeployVM(context.Background(), args.ID, "", args.Password, args.PublicKey, false, "")
	})

	type InstanceSecurityGroupOptions struct {
		ID            string `help:"Instance ID"`
		SecurityGroup string `help:"Security Group ID or Name"`
	}

	shellutils.R(&InstanceSecurityGroupOptions{}, "instance-set-secgrp", "Attach a disk to intance", func(cli *azure.SRegion, args *InstanceSecurityGroupOptions) error {
		return cli.SetSecurityGroup(args.ID, args.SecurityGroup)
	})
}
