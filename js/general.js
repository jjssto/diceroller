import {    
    hasFocus,
    highlightOwnRolls,
    addCol,
    init,
    initReset
} from "./fun.js";

init(createRow)
initReset(reset)

function createRow(drow, id) {
    var row = document.createElement("tr")
    const dat = drow[0]
    const isOwnRoll = drow[1]
    if (isOwnRoll == "1" && highlightOwnRolls) {
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

function setDice() {
    var sel = document.querySelector("#s_d20");
    const d20 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d12");  
    const d12 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d10");  
    const d10 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d8");
    const d8 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d6");
    const d6 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d4");
    const d4 = sel.options[sel.selectedIndex].value; 
    
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
    document.querySelector("#s_d20").value = 0;
    document.querySelector("#s_d12").value = 0;
    document.querySelector("#s_d10").value = 0;
    document.querySelector("#s_d8").value = 0;
    document.querySelector("#s_d6").value = 0;
    document.querySelector("#s_d4").value = 0;
 }

document.getElementById("f_roll").addEventListener("submit", event => {
    event.preventDefault()
    const loc = location.href
    const player_id = "0"
    const dice = setDice(); 
    const chr = document.getElementById("f_name").value;
    const action = document.getElementById("f_action").value;

    fetch(loc, {
        method: "POST",
        headers: {
            "contentType": "application/json"
        },
        body: JSON.stringify({
            "dice": dice,
            "char": chr,
            "action": action
        })
    })
})