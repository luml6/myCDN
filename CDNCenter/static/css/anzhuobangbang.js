function abbIndexOf(a) {
	return (location.href.indexOf(a) > -1);
}
function getQueryString(a)
{
		var reg = new RegExp("(^|&)"+ a +"=([^&]*)(&|$)");
    	var r = window.location.search.substr(1).match(reg);
     	if(r!=null)return  unescape(r[2]); return null;
}
function bindingData(a){
	var aThead = $("thead th");
							var aTQ = "";
							var aT = ","
							var aA = new Array();
							var aAH = new Array();
							for (var n = 0; n < aThead.length; n++) {
								if (aThead.length == (n + 1)) {
									aT = "";
								}
								aA.push(aThead[n].id);
								aAH.push(aThead[n].innerHTML);
								if (aThead[n].innerHTML.indexOf("{0}") > 0) {
									aThead[n].innerHTML = ""
								}
								aTQ += aThead[n].id + aT
							}
							var aH = "";
							var aPage = "";
							$.ajax({
								url: "a?page="+getQueryString("p"),
								type: "GET",
								context: document.body,
								success: function(result) {
									//result = eval("(" + result + ")"); /*修改*/
						
									var aResult1 = result[1];
									var aResult0 = result[0];
									
						
									for (var j = 1; j < (aResult0["totalCount"]+1); j++) {
										aPage += '<a href="channel.html?page='+j+'" class="item">'+j+'</a>'
									}
						
									$("#aPage").html(aPage);
						
									for (var n = 0; n < aResult1.length; n++) {
										aH += "<tr>";
										for (var m = 0; m < aA.length; m++) {
											if (aAH[m].indexOf("{0}") > 0) {
												aH += "<td>" + aAH[m].replace("{0}", aResult1[n][aA[m]]) + "</td>";
											} else {
												aH += "<td>" + aResult1[n][aA[m]] + "</td>";
											}
										}
										aH += "</tr>";
									}
									$("tbody").append(aH);
								}
							});
	
	
}

						
						
$(document).ready(function() {
	$('.ui.fluid.search.dropdown').dropdown();
	$('.ui.menu .ui.dropdown').dropdown({
		on: 'hover'
	});
	$("#menu").html('<a class="item" href="slavemanager.html" id="slavemanager">节点管理 </a><a class="item" href="download.html" id="download">文件下载情况</a>');

	if (abbIndexOf("slavemanager")) {
		$("#slavemanager").addClass("active");
	}
	if (abbIndexOf("download")) {
		$("#download").addClass("active");
	}
	$("#nav").html(' <img src="static/img/minic.png" class="image" style="margin:15px auto 0px auto"></br>')
	

});