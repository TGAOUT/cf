package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alirds"
)

var (
	rdslsFlushCache            bool
	rdslsRegion                string
	rdslsEngine                string
	rdslsSpecifiedDBInstanceID string
)

func init() {
	RootCmd.AddCommand(rdsCmd)
	rdsCmd.AddCommand(rdslsCmd)
	rdsCmd.PersistentFlags().BoolVar(&rdslsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
	rdslsCmd.Flags().StringVarP(&rdslsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	rdslsCmd.Flags().StringVarP(&rdslsSpecifiedDBInstanceID, "DBInstanceID", "i", "all", "指定数据库实例 ID (Set DBInstance ID)")
	rdslsCmd.Flags().StringVarP(&rdslsEngine, "engine", "e", "all", "指定数据库类型 (Set DBInstance Type)")
}

var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "执行与云数据库相关的操作 (Perform rds-related operations)",
	Long:  "执行与云数据库相关的操作 (Perform rds-related operations)",
}

var rdslsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的云数据库 (List all DBInstances)",
	Long:  "列出所有的云数据库 (List all DBInstances)",
	Run: func(cmd *cobra.Command, args []string) {
		alirds.PrintDBInstancesList(rdslsRegion, running, rdslsSpecifiedDBInstanceID, rdslsEngine, rdslsFlushCache)
	},
}
