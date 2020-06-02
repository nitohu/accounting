$(document).ready(function() {
    let btns = document.getElementsByClassName("deleteEntry")
    let tbody = document.getElementById("transaction_list")
    const currency = $(tbody).attr("data-currency")

    getAccounts()

    function addTableButtonEvent() {
        for (let i = 0; i < btns.length; i++) {
            let btn = btns[i]
            btn.addEventListener("click", function(e) {
                let id = $(this).attr("data-id")
                deleteTransaction(id)
            })
        }
    }

    function generateTableItems(data) {
        for(let index in data) {
            let item = data[index]
            let row = document.createElement("tr")

            let child = document.createElement("a")
            child.setAttribute("href", "/transactions/form?id="+item.ID)
            child.innerText = item.Name
            let parent = document.createElement("td")
            parent.appendChild(child)
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.Amount.toFixed(2) + " " + currency
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.FromAccountName
            row.appendChild(parent)

            parent = document.createElement("td")
            parent.innerText = item.Category.Name
            parent.setAttribute("style", "color:"+item.Category.Hex)
            row.appendChild(parent)
            
            parent = document.createElement("td")
            parent.innerText = item.TransactionDateStr
            row.appendChild(parent)
            
            parent = document.createElement("td")
            parent.innerText = item.ToAccountName
            row.appendChild(parent)

            child = document.createElement("i")
            child.setAttribute("data-id", item.ID)
            child.setAttribute("class", "material-icons")
            child.innerText = "X"
            parent = document.createElement("td")
            parent.setAttribute("data-id", item.ID)
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

        xhr.open("GET", "/api/transactions", true)
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

    function deleteTransaction(id) {
        let xhr = new XMLHttpRequest()
        let data = JSON.stringify({"ID": Number.parseInt(id)})

        xhr.open("DELETE", "/api/transactions/delete", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                let m = this.responseText.replace(/'/g, '"')
                let msg = JSON.parse(m)
                console.log(msg.success)
                getAccounts()
            } else if (this.readyState == 4 && (this.status == 400 || this.status == 405)) {
                let m = this.responseText.replace(/'/g, '"')
                let msg = JSON.parse(m)
                console.warn(msg.error)
                getAccounts()
            }
        }
        xhr.send(data)
    }
})