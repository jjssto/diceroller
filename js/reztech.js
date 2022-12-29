import {    
    hasFocus,
    addCol,
    init,
    initReset
} from "./fun.js";

init(createRow);
initReset(reset);

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

function setDice() {
    var sel = document.querySelector("#s_attribute");
    const attr = sel.options[sel.selectedIndex].value;
    sel = document.querySelector("#s_skill");
    const skill = sel.options[sel.selectedIndex].value;
    const onlyAttr = document.querySelector("#i_attribute_only").checked;
   
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
    } else if (skill == 0) {
        for (let i = 0; i < attr; i++) {
            if (first) {
                first = false
            } else {
                text += ", "
            }
            text += "6"
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
    document.getElementById("s_attribute").value = 0;
    document.getElementById("s_skill").value = 0;
    document.getElementById("i_attribute_only").checked = false;
  }


document.getElementById("f_roll").addEventListener("submit", event => {
    event.preventDefault()
    const loc = location.href
    const player_id = "0"
    const dice = setDice() 
    const name = document.getElementById("f_name").value;
    const action = document.getElementById("f_action").value;
   
    fetch(loc, {
        method: "POST",
        headers: {
            "contentType": "application/json"
        },
        body: JSON.stringify({
            "dice": dice,
            "char": name,
            "action": action
        })
    })
});

