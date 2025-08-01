USER_LOGIN = {

    id: "",
    passphrase: ""
}

COMMENT = {

    title: "",
    content: ""
}

URL = ""
APPKEY = ""

IMAGE_TITLE = ""
IMAGE_GROOM = ""
IMAGE_BRIDE = ""

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



    storyList.innerHTML = ""

    for(let i = 0; i < entryLen; i++){

        let oe = orderedEntry[i]

        let content = `

                <div class="row-story" onclick="window.open('/story/r/${oe.id}','_self')">
                    <div class="col-md-7 col-sm-12">
                        <div class="my-3">
                            <div class="h4">${oe.title}</div>
                            <ul>
                                <li>
                                    <i class="text-muted fas fa-map-marker-alt"></i>
                                    <span class="pl-2 text-muted">${oe.intro}</span>
                                </li>
                                <li class="pt-2">
                                    <i class="text-muted far fa-calendar-alt"></i>
                                    <span class="pl-2 text-muted">${oe.dateMarked}</span>
                                </li>
                                <li class="pb-2">
                                    <i class="text-muted far fa-calendar-alt"></i>
                                    <span class="pl-2 text-muted">(클릭하면 스토리가 보입니다)</span>
                                </li>
                            </ul>
                        </div>
                    </div>
                    <div class="col-md-5 col-sm-12">
                        <div class="my-3">
                            <img alt="Wedding Party" class="img-fluid" src="${oe.primaryMediaName}"/>
                        </div>
                    </div>
                </div>


        `

        storyList.innerHTML += content


    }
  
}



async function getImageList(){


    let resp = await fetch("/api/image/list", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get comment list")

        return

    }

    let imageList = document.getElementById("image-rows")

    let imageEntry = JSON.parse(result.reply)

    let entrylen = imageEntry.length 

    imageList.innerHTML = ""

    for(let i = 0 ; i < entrylen; i++){

        if(i == 0){
            continue
        }

        let ie = imageEntry[i]

        if (IMAGE_GROOM == ""){
            IMAGE_GROOM = ie.name
            let g = document.getElementById("couple-groom")
            g.innerHTML += `
            <img alt="Groom" class="img-fluid" src="/${IMAGE_GROOM}"/>
            `
            continue
        } else if (IMAGE_BRIDE == ""){
            IMAGE_BRIDE = ie.name
            let b = document.getElementById("couple-bride")
            b.innerHTML += `
            <img alt="Bride" class="img-fluid" src="/${IMAGE_BRIDE}"/>
            `
            continue
        }

        let ieEl = `
            <div class="card" data-groups="[&quot;party&quot;,&quot;wedding&quot;]">
                <a data-gallery="ww-gallery" data-toggle="lightbox">
                    <img alt="Gallery Pic 2" class="img-fluid" src="/${ie.name}"/>
                </a>
            </div>
        `

        imageList.innerHTML += ieEl
    }
}


async function registerComment(){


    let c_title = document.getElementById("comment-title").value 

    if(c_title == ""){
  
        alert("방명록 남기시는 분이 누구인지 알려주세요~")
    
        return
    
    }


    let c_content = document.getElementById("comment-content").value 


    if(c_content == ""){
  
        alert("방명록 내용을 채워주세요~")
    
        return
    
    }


    let com = JSON.parse(JSON.stringify(COMMENT))

    com.title = c_title
    com.content = c_content

    let req = {
        data: JSON.stringify(com)
    }

    let commentSection = document.getElementById("comment-section")

    commentSection.innerHTML = `
    
    <img src="/public/loading.gif"/>    

    `



    let resp = await fetch(`/api/comment/register`, {
        body: JSON.stringify(req),
        method: "POST"
    })


    let result = await resp.json()

    if(result.status != "success"){

        alert("방명록 남기기에 실패 했습니다, 다시 시도해주세요: ", + result.reply)

        location.href = "/"

        return
    }


    alert("방명록을 성공적으로 남겼습니다: " + result.reply + "\n확인 후 게시 예정입니다.\n감사합니다 ^^")

    location.href = "/"

}

async function registerCommentSudo(){



    let cmd = document.getElementById("comment-sudo-cmd").value 

    if(cmd == ""){
  
        alert("specify cmd")
    
        return
    
    }

    let file = document.getElementById("comment-sudo-file")

    const form = new FormData()

    form.append("file", file.files[0]);

    if(cmd == "allow"){
        let resp = await fetch("/api/comment/sudo/allow", {
                body: form,
                method: "POST"
        })

        console.log(resp)

        let resp_json = await resp.json()

        if(resp_json.status != "success"){
            alert("failed to handle allow command")
            return
        } 

        alert("successfully handled allow command: \n" + resp_json.reply)

    } else if (cmd == "block") {
        let resp = await fetch("/api/comment/sudo/block", {
                body: form,
                method: "POST"
        })

        let resp_json = await resp.json()

        if(resp_json.status != "success"){
            alert("failed to handle block command")
            return
        }

        alert("successfully handled block command: \n" + resp_json.reply)

    } else {
        alert("invalid cmd: " + cmd)
        return
    }


}


async function getCommentList(){

    if(GOOGLE_COMMENT_Y == ""){
        return
    }

    let gce = document.getElementById("google-comment")

    gce.innerHTML = GOOGLE_COMMENT_EL

    let resp = await fetch("/api/comment/list", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get comment list")

        return

    }

    let commentList = document.getElementById("comment-rows")

    let commentEntry = JSON.parse(result.reply)

    let md = new Remarkable();
    let entrylen = commentEntry.length 

    let rawText = ""
    let mdRenderd = ""



    for(let i = 0 ; i < entrylen; i++){

        let ce = commentEntry[i]

        rawText += "### " + ce.title + "\n"
        rawText += ce.content + "\n\n"

    }

    mdRenderd = md.render(rawText)

    commentList.innerHTML = mdRenderd

}

async function getAppShare(){

    let kse = document.getElementById("gift-shares")

    let innerhtml = SHARE_EL

    if(KAKAO_SHARE_Y != ""){
        
        innerhtml += KAKAO_SHARE_EL
    }

    kse.innerHTML = innerhtml;
    if(KAKAO_SHARE_Y == ""){
        
        return
    }

    let resp = await fetch("/api/appkey", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get appkey")

        return

    }

    titleAppkey = result.reply 

    ta = titleAppkey.split(":")

    let title = ta[0]

    APPKEY = ta[1]

    Kakao.init(APPKEY);

    thisUrl = window.location.href

    Kakao.Share.createDefaultButton({
        container: '#kakaotalk-sharing-btn',
        objectType: 'feed',
        content: {
          title: title,
          imageUrl:
            `${thisUrl}${IMAGE_TITLE}`,
          link: {
            mobileWebUrl: thisUrl,
            webUrl: thisUrl,
          },
        },
        buttons: [
          {
            title: '소식 보기',
            link: {
              mobileWebUrl: thisUrl,
              webUrl: thisUrl,
            },
          },
        ],
      });
}


function copyUrlToClipboard(){


    thisUrl = window.location.href
    
    navigator.clipboard.writeText(thisUrl);

    alert("url이 클립보드에 복사되었습니다")
}


async function getGiftPage(){


    let resp = await fetch("/api/gift", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get gift page")

        return

    }


    window.open(result.reply, '_blank').focus();

}
