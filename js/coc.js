import "./jquery-3.6.1.js"

function addCol(row, text) {
    col = document.createElement("td")
    col.innerHtml = text 
    row.appendChild(col)
}

function createRow(drow) {
    var row = document.createElement("tr")
    addCol(row, drow[0])
    const dat = drow[1]
    addCol(row, dat.T)
    addCol(row, dat.P)
    addCol(row, dat.D)
    addCol(row, dat.R)
    return row
}

function insertRolls(data) {
    const tbody = document.getElementById("tbody_rolls")
    const first_row = tbody.firstChild
    var roll_id
    if (first_row) {
        roll_id = first_row.getElementById("roll_id").val()
    } else {
        roll_id = 0
    }

    for (i=0; i<length(data); i++) {
        var drow = data[i]
        if (roll_id >= drow[0]) {
            continue
        }
        localStorage.setItem("last_roll", drow[0])
        const row = createRow(drow)
        if (roll_id > 0) { 
            first_row.insertBefore(row)
        } else {
            tbody.appendChild(row)
        }
    }
}

function getRolls() {
    var target = location.href.replace("room", "rolls")
    const last_roll = localStorage.getItem("last_roll")
    if (last_roll != "") {
        target += last_roll
    }
   $.ajax({
        url: target,
        method: "GET",
        dataType: "app/json",
        success: insertRolls  
   }) 
}

$(document).ready(function(){
    localStorage.setItem("last_roll", "")
    getRolls()


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
            method: "GET",
            data: data,
            contentType: "app/json",
            dataType: "app/json",
            success: insertRolls
        })
    })
})

window.setInterval(getRolls, 1000)