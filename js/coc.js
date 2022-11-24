import "./jquery-3.6.1.js"

$(document).ready(function(){
    $("#f_roll").submit((event) => {
        event.preventDefault()
        const loc = location.href
        const player_id = "0"
        const mod = $("input[name='mod']:checked").val()
        var data = "{"
        data += '"player": "' + player_id + '",'
        data += '"mod": "' + mod + '",'
        data += '"char": "' + $("#f_name").val() + '",'
        data += '"action": "' + $("#f_action").val() + '"'
        data += "}"
        $.ajax({
            url: loc,
            method: "POST",
            data: data,
            contentType: "app/json",
        })
    })
})