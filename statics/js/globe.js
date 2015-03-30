require([
	"dojo/ready",
	"dojo/parser",
	"dojo/dom-style",
	"dojo/_base/fx",
	"dojo/_base/lang",
	"dijit/registry",
	"dijit/layout/ContentPane",
	"dojo/when",
	"dojo/io/script",
	"dojo/request"
], function(ready,parser,domStyle,fx,lang,registry,ContentPane,when,ioScript,request){


ready(function(){
	loadPage();
	index.showPage("站点1","/html/register.html","/js/qq.js");
});


/**
 * 全局方法，整个页面任意位置都能调用
 * tabName: tab显示的title名称
 * htmlUrl: 加载的html页面路径
 * jsUrl：加载的js文件路径，可以为空
 * cssUrl：加载的css样式文件路径，可以为空
 **/
lang.setObject("index.showPage",function(tabName,htmlUrl,jsUrl){
	// alert(1);
	var content = registry.byId("content");
	content.setHref( htmlUrl);
	ioScript.get({url:jsUrl});
	// // alert(workspaces);
	// content.destroy();

	// new ContentPane({
 //      href:htmlUrl
 //    }, "workspaces").refresh();


// var myCp= registry.byId("content");
// myCp.set("content", function(){
//     console.log("Download complete!");
// });
// myCp.set("href", htmlUrl);



// dojo.create("div", { style: { color: "red" },href:htmlUrl,"data-dojo-type":"dijit/layout/ContentPane",id:"content" }, "workspaces");

	// alert(3);
	
	// content.refresh()
	// //tab的id直接用htmlUrl路径命名把
	// var tabOne = registry.byId(htmlUrl);
	// if(typeof tabOne === "undefined"){
	// 	tabOne = new ContentPane({
	// 		id:htmlUrl,
	// 		title:tabName,
	// 		closable:true,  //是否显示关闭按钮
	// 		//content: "We are known for our drinks."
	// 		href:htmlUrl
	// 	});
	// 	workspaces.addChild(tabOne);
	// 	workspaces.selectChild(tabOne);
	// 	ioScript.get({url:jsUrl});
	// 	// when(tabContainer.addChild(tabOne),
	// 	// 	when(tabContainer.selectChild(tabOne),function(){
	// 	// 		if(jsUrl != null && jsUrl != ""){
	// 	// 			ioScript.get({url:jsUrl});
	// 	// 		}
	// 	// 	})
	// 	// );
	// }else{
	// 	//选中tab
	// 	workspaces.selectChild(tabOne);
	// }
});

lang.setObject("index.ajax" , function(url , data){
	this.then = function(callback, errback, progback){
		request(url, data).then(function(data){
			callback(data);
		} , function(err){
			errback(err);
		} , function(e){
			progback(e);
		});
	};
	return this;
});



function loadPage(){
	parser.parse().then(function(objects){
		fx.fadeOut({  //Get rid of the loader once parsing is done
			node: "index_loader",
			duration: 0,
			onEnd: function() {
				domStyle.set("index_loader","display","none");
			}
		}).play();
	});
}


});