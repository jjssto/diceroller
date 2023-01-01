import {    
    hasFocus,
    addCol,
    init,
    highlightOwnRolls
} from "./fun.js";

init(createRow);

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


document.getElementById("f_roll").addEventListener("submit", event => {
    event.preventDefault()
    const loc = location.href
    const player_id = "0"
    const checked = document.querySelector('input[name="mod"]:checked');
    var mod = "";
    if (checked != null) {
        mod = checked.value;
    } else {
        mod = "0";
    }
    const chr = document.getElementById("f_name").value;
    const action = document.getElementById("f_action").value;
    fetch(loc, {
        method: "POST",
        headers: {
            "contentType": "application/json"
        },
        body: JSON.stringify({
            "mod": mod,
            "char": chr,
            "action": action 
        })
    })
})