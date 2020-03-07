$("document").ready(function() {
    var categoryID = 0
    var tbody = document.getElementById("categoryTable")

    getCategories()
    let addLineBtn = document.getElementById("addLineBtn")
    addLineBtn.addEventListener("click", function() {
        categoryID = 0
    })

    let formBtn = document.getElementById("formBtn")
    formBtn.addEventListener("click", saveCategories)

    function debounce(func, wait, immediate) {
        var timeout;

        return function executedFunction() {
            var context = this;
            var args = arguments;
                
            var later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
            };

            var callNow = immediate && !timeout;
            
            clearTimeout(timeout);

            timeout = setTimeout(later, wait);
            
            if (callNow) func.apply(context, args);
        };
    };

    function saveCategories() {
        let name = document.getElementById("name").value
        let hex  = document.getElementById("hex").value

        let data = {
            "ID": categoryID,
            "Name": String(name),
            "Hex": String(hex),
        }

        let xhr = new XMLHttpRequest()
        xhr.open("POST", "/api/categories/create")
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = debounce(getCategories, 200, false)
        xhr.onload = () => {
            console.log(xhr.status)
            console.log(xhr.response)
        }
        xhr.send(JSON.stringify(data))
        $("#formModal").modal('hide')
    }

    function getCategoryInformation(id) {
        let xhr = new XMLHttpRequest()
        xhr.open("POST", "/api/categories", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                res = JSON.parse(this.responseText)
                document.getElementById("name").value = res["Name"]
                document.getElementById("hex").value = res["Hex"]
            } else if (this.readyState == 4) {
                console.error(JSON.parse(this.responseText))
            }
        }
        d = {"ID": Number.parseInt(id)}
        xhr.send(JSON.stringify(d))
    }

    function getCategories() {
        let data = {
            "ID": 0,
        }
        let xhr = new XMLHttpRequest()
        xhr.open("GET", "/api/categories", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                // Remove current categories in DOM
                while (tbody.lastElementChild) {
                    tbody.removeChild(tbody.lastElementChild)
                }
                // Append new categories
                let categories = JSON.parse(this.responseText)
                for(let i in categories) {
                    let category = categories[i]
                    let content = document.createElement("tr")
                    content.setAttribute("style", 'background-color: '+category["Hex"])
                    content.setAttribute("class", "categoryItem")
                    content.setAttribute("id", category["ID"])

                    let td = document.createElement("td")
                    td.innerHTML = category["Hex"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    let a = document.createElement("a")
                    a.setAttribute("href", "#")
                    a.setAttribute("data-toggle", "modal")
                    a.setAttribute("data-target", "#formModal")
                    a.innerHTML=category["Name"]
                    a.addEventListener("click", function() {
                        let id = this.parentElement.parentElement.id
                        c = getCategoryInformation(id)
                    })
                    td.appendChild(a)
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["CreateDate"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["LastUpdate"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["TransactionCount"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.setAttribute("class", "deleteEntry")
                    a = document.createElement("i")
                    a.setAttribute("class", "zmdi zmdi-close")
                    a.addEventListener("click", deleteCategory)
                    td.appendChild(a)
                    content.appendChild(td)
                    
                    // let row = document.getElementById("addLineRow")
                    // tbody.insertBefore(content, row)
                    tbody.appendChild(content)
                }
            }
        }
        xhr.send(JSON.stringify(data))
    }

    function deleteCategory(e) {
        data = {
            "ID": Number.parseInt(e.path[2].id)
        }
        tbody.removeChild(e.path[2])
        let xhr = new XMLHttpRequest()
        xhr.open("DELETE", "/api/categories/delete", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 400) {
                console.error(this.response)
            }
        }
        console.log(JSON.stringify(data))
        xhr.send(JSON.stringify(data))
        console.log(xhr)
    }
})