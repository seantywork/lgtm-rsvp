USER_LOGIN = {

    id: "",
    passphrase: ""
}

async function userSignin(){



    let u_id = document.getElementById("user-id").value 

    if(u_id == ""){
  
        alert("no user id")
    
        return
    
    }


    let u_pw = document.getElementById("user-pw").value 


    if(u_pw == ""){
  
        alert("no user pw")
    
        return
    
    }



    let uc = JSON.parse(JSON.stringify(USER_LOGIN))

    uc.id = u_id
    uc.passphrase = u_pw

    let req = {
        data: JSON.stringify(uc)
    }

    let resp = await fetch(`/api/signin`, {
        body: JSON.stringify(req),
        method: "POST"
    })


    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to login")

        return
    }

    alert("successfully logged in: " + result.reply)

    location.href = "/"


}




async function getArticleList(){


      let resp = await fetch("/api/story/list", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get sample list")

        return

    }


    let storyList = document.getElementById("story-columns")

    let contentEntry = JSON.parse(result.reply)
    

    let entryLen = contentEntry.length

    let orderedEntry = []

    for(let i = 0 ; i < entryLen; i++){

        let len = contentEntry.length

        let idx = 0
        let num = 999999

        for(let j = 0; j < len; j++){

            let entry = contentEntry[j]

            let dateMarked = entry.dateMarked

            let dateStr = dateMarked.replace("-", "")

            let dateNum = parseInt(dateStr)

            if(dateNum < num){
                idx = j
                num = dateNum
            }
        }

        let newEntry = []

        for(let j = 0; j < len; j++){

            if(j == idx){
                orderedEntry.push(contentEntry[j])
                continue
            }

            newEntry.push(contentEntry[j])
        }

        contentEntry = newEntry

    }

    console.log(orderedEntry)

    storyList.innerHTML = ""

    for(let i = 0; i < entryLen; i++){

        let oe = orderedEntry[i]

        let content = `
            <div class="gift">
                <img alt="와인셀러" class="img-fluid gift-img gift-selected" src="${oe.primaryMediaName}"/>
                <div class="btn-group gift-btn-group" role="group">
                    <button class="btn btn-default gift-btn"
                            onclick="window.open('/story/r/${oe.id}')">
                        ${oe.dateMarked}</button>

                </div>
            </div>        
        
        `

        storyList.innerHTML += content


    }
  
}

/*
                   <button class="btn btn-default gift-btn"
                            onclick="window.open('/story/r/${oe.id}')">
                        <i class="fa fa-search"></i></button>                    
<button class="btn btn-default gift-btn gift-send" data-name="와인셀러"
                            onclick="alert('다른분에게 예약된 선물입니다.');"><i
                            class="fa fa-gift"></i></button>
*/