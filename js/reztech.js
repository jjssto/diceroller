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
    var d12 = ""
    var d8 = ""
    var d6 = ""
    for(let i in dat.D) {
        switch (dat.D[i].E) {
            case "12":
                d12 += dat.D[i].R 
                d12 += ", "
                break;
            case "8":
                d8 += dat.D[i].R 
                d8 += ", "
                break;
            case "6":
                d6 += dat.D[i].R 
                d6 += ", "
                break;
        }
    }
    addCol(row, d12)
    addCol(row, d8)
    addCol(row, d6)
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
    const attr = $("#s_attribute").find(":selected").val()
    const skill = $("#s_skill").find(":selected").val()
    const onlyAttr = $("#i_attribute_only").prop("checked")
    
    var first = true
    var text= "["
    if (attr <= 0 && skill <= 0) {
        // do nothing
    } else if (onlyAttr) {
         for (let i = 0; i < attr; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "8"
        }
    } else if (attr >= skill) {
        for (let i = 0; i < skill; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "12"
        }
        for (let i = 0; i < attr - skill; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "8"
        }
    } else {
         for (let i = 0; i < attr; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "12"
        }
        for (let i = 0; i < skill - attr; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "6"
        }
    }
    text += "]"
    return text
}

function reset() {
    $("#s_attribute").find("option.o1").prop("selected", true) 
    $("#s_skill").find("option.o1").prop("selected", true) 
    $("#i_attribute_only").prop("checked", false)
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
    
    $("#b_reset").click(() => {
        reset()
    })
})

window.setInterval(getRolls, 1000)