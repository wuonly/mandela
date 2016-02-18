/*
	加载本地配置文件中的超级节点地址
*/

package addr_manager

func init() {
	registerFunc(loadByLocal)
}

/*
	通过自定义目录服务器获得超级节点地址
*/
func loadByLocal() {

}
