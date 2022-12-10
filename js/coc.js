import "./jquery-3.6.1.js"

function addCol(row, text) {
    const col = row.insertCell()
    col.appendChild(document.createTextNode(text))
}

function createRow(drow, id) {
    var row = document.createElement("tr")
    const dat = drow[0]
    const isOwnRoll = drow[1]
    if (isOwnRoll == "1") {
        row.className = "my_roll"
    }
    addCol(row, id)
    addCol(row, dat.T)
    const colP = row.insertCell()
    colP.appendChild(document.createTextNode(dat.P))
    colP.className = "my_name"
    addCol(row, dat.A)
    const colD = row.insertCell()
    var text = ""
    for(let i in dat.D) {
        text += dat.D[i].R 
        text += ", "
    }
    colD.appendChild(document.createTextNode(text))
    addCol(row, dat.R)
    return row
}

function insertRolls(data_raw) {
    const data = JSON.parse(data_raw)
    const tbody = document.getElementById("tbody_rolls")
    const first_row = tbody.firstChild
    var roll_id
    if (first_row) {
        roll_id = first_row.firstChild.firstChild.textContent
    } else {
        roll_id = -1
    }
    for (let i in data) {
        const drow = data[i]
        localStorage.setItem("last_roll", i)
        const row = createRow(drow, i)
        if (i > roll_id) { 
            //if (roll_id >= 0) {
            //    tbody.insertBefore(row, first_row)
            //} else {
                tbody.appendChild(row)
                row.scrollIntoView(true)
            //}
        }
    }
}

function getRolls() {
    var target = location.href.replace("room/", "rolls/")
    const last_roll = localStorage.getItem("last_roll")
    if (last_roll != "") {
        target += "/" + last_roll
    }
   $.ajax({
        url: target,
        method: "GET",
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

window.setInterval(getRolls, 1000)