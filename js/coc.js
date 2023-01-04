// This code is licensed under the MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


import {    
    hasFocus,
    addCol,
    formatTime,
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
    addCol(row, formatTime(dat.T))
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