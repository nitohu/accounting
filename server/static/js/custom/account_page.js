$(document).ready(function() {
    getAccounts()

    let btns = document.getElementsByClassName("deleteEntry")
    let tbody = document.getElementById("account_list")

    function addTableButtonEvent() {
        for (let i = 0; i < btns.length; i++) {
            let btn = btns[i]
            btn.addEventListener("click", deleteAccount)
        }
    }

    function generateTableItems(data) {
        for(let index in data) {
            let item = data[index]
            let row = document.createElement("tr")

            let child = document.createElement("a")
            child.setAttribute("href", "/accounts/form?id="+item.ID)
            child.innerText = item.Name
            let parent = document.createElement("td")
            parent.appendChild(child)
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.Balance.toFixed(2)
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.BankName
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.Iban
            row.appendChild(parent)

            child = document.createElement("i")
            child.setAttribute("account-id", item.ID)
            child.setAttribute("class", "material-icons")
            child.innerText = "X"
            parent = document.createElement("td")
            parent.setAttribute("account-id", item.ID)
            parent.setAttribute("class", "deleteEntry")
            parent.appendChild(child)
            row.appendChild(parent)

            tbody.appendChild(row)
        }
        addTableButtonEvent()
    }

    function getAccounts() {
        let data = {"ID": 0}
        let xhr = new XMLHttpRequest()

        xhr.open("GET", "/api/accounts", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                // Remove all list items from the list
                while(tbody.lastElementChild) {
                    tbody.removeChild(tbody.lastElementChild)
                }
                
                // Generate & append new list items to the list
                let data = JSON.parse(this.responseText)
                generateTableItems(data)
            }
        }
        xhr.send(data)
    }

    function deleteAccount(e) {
        let id = $(e.target).attr("account-id")
        data = {"ID": Number.parseInt(id)}
        
        let xhr = new XMLHttpRequest()
        xhr.open("DELETE", "/api/accounts/delete", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4) {
                let m = this.response.replace(/'/g, '"')
                let msg = JSON.parse(m)
                if (this.status == 400) {
                    console.error(msg.error)
                } else if (this.status == 200) {
                    getAccounts()
                    console.log(msg.success)
                } else if (this.status == 403) {
                    console.warn(msg.error)
                } else {
                    console.error(this.response)
                }
            }
        }
        xhr.send(JSON.stringify(data))
    }
})
