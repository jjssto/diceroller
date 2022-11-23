import "./jquery-3.6.1.js"
$(document).ready(function(){
        $("button").click(function(){
           $.post("/",
           {
            id: $(this).val()
           },
           function(data, status){
               json = JSON.parse(data)
               document.location.href = json.id 
           }
           ) 
        })
    })
