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
    highlightOwnRolls,
    formatTime,
    addCol,
    init
} from "./fun.js";

init(createRow, rollDice);

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
    const attr = document.querySelector("#s_attribute").value;
    const skill = document.querySelector("#s_skill").value;
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


function rollDice() {
    const loc = location.href
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
};

