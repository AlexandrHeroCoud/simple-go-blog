$(function(){
    $("#content").bind("input change", function(){
        $.post("/gethtml", {md: $("#content").val()}, function(response){
                $("#markdown_html").html(response.html)
            })
    })})