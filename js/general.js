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
    var d20 = ""
    var d12 = ""
    var d10 = ""
    var d8 = ""
    var d6 = ""
    var d4 = ""
    for(let i in dat.D) {
        switch (dat.D[i].E) {
            case "20":
                d20 += dat.D[i].R 
                d20 += ", "
                break;
            case "12":
                d12 += dat.D[i].R 
                d12 += ", "
                break;
            case "10":
                d10 += dat.D[i].R 
                d10 += ", "
                break;
            case "8":
                d8 += dat.D[i].R 
                d8 += ", "
                break;
            case "6":
                d6 += dat.D[i].R 
                d6 += ", "
                break;
            case "4":
                d4 += dat.D[i].R 
                d4 += ", "
                break;
        }
    }
    addCol(row, d20)
    addCol(row, d12)
    addCol(row, d10)
    addCol(row, d8)
    addCol(row, d6)
    addCol(row, d4)
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
            tbody.appendChild(row)
           row.scrollIntoView(true)
        }
    }
}

function getRolls() {
    var target = location.href.replace("room", "rolls")
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

function setDice() {
    const d20 = $("#s_d20").find(":selected").val()
    const d12 = $("#s_d12").find(":selected").val()
    const d10 = $("#s_d10").find(":selected").val()
    const d8 = $("#s_d8").find(":selected").val()
    const d6 = $("#s_d6").find(":selected").val()
    const d4 = $("#s_d4").find(":selected").val()
    
    var first = true
    var isEmpty = false
    var text= "["
    for (let i = 0; i < d20; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "20"
        isEmpty = false
    }
    for (let i = 0; i < d12; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "12"
        isEmpty = false
    }
    for (let i = 0; i < d10; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "10"
        isEmpty = false
    }
    for (let i = 0; i < d8; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "8"
        isEmpty = false
    }
    for (let i = 0; i < d6; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "6"
        isEmpty = false
    }
    for (let i = 0; i < d4; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "4"
        isEmpty = false
    }
    text += "]"
    return text
}

function reset() {
    $("#s_d20").find("option.o1").prop("selected", true) 
    $("#s_d12").find("option.o1").prop("selected", true) 
    $("#s_d10").find("option.o1").prop("selected", true) 
    $("#s_d8").find("option.o1").prop("selected", true) 
    $("#s_d6").find("option.o1").prop("selected", true) 
    $("#s_d4").find("option.o1").prop("selected", true) 
}

$(document).ready(function(){
    localStorage.setItem("last_roll", "")
    getRolls()


    $("#f_roll").submit((event) => {
        event.preventDefault()
        const loc = location.href
        const player_id = "0"
        const dice = setDice() 
        var data = "{"
        data += '"dice": "' + dice + '",'
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
    
    $("#b_reset").click((event) => {
        reset()
    })
})

window.setInterval(getRolls, 1000)