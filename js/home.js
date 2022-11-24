import "./jquery-3.6.1.js"
$(document).ready(function(){
    $("button").click(function(){
        $.post("/",
           {
            id: $(this).val()
           },
           function(data, status){
               if (status == "success") {
                const loc = `${window.location.href}room/${data}`
                window.location.assign(loc) 
               }
           }
           ) 
        })
    })
