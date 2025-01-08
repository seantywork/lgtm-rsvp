
ARTICLE_INFO = {
    "title": "",
    "content": "",
    "dateMarked": "",
    "primaryMediaName":""

}


const ImageTool = window.ImageTool;



async function mediaUploader(fileData){

  const form = new FormData()

  form.append("file", fileData)


  let resp = await fetch("/api/media/upload", {
        body: form,
        method: "POST"
  })

  let resp_json = await resp.json()

  return resp_json

}


/**
 * Module to compose output JSON preview
 */
const cPreview = (function (module) {
  /**
   * Shows JSON in pretty preview
   * @param {object} output - what to show
   * @param {Element} holder - where to show
   */
  module.show = function(output, holder) {
    /** Make JSON pretty */
    output = JSON.stringify( output, null, 4 );
    /** Encode HTML entities */
    output = encodeHTMLEntities( output );
    /** Stylize! */
    output = stylize( output );
    holder.innerHTML = output;
  };

  /**
   * Converts '>', '<', '&' symbols to entities
   */
  function encodeHTMLEntities(string) {
    return string.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  /**
   * Some styling magic
   */
  function stylize(string) {
    /** Stylize JSON keys */
    string = string.replace( /"(\w+)"\s?:/g, '"<span class=sc_key>$1</span>" :');
    /** Stylize tool names */
    string = string.replace( /"(paragraph|quote|list|header|link|code|image|delimiter|raw|checklist|table|embed|warning)"/g, '"<span class=sc_toolname>$1</span>"');
    /** Stylize HTML tags */
    string = string.replace( /(&lt;[\/a-z]+(&gt;)?)/gi, '<span class=sc_tag>$1</span>' );
    /** Stylize strings */
    string = string.replace( /"([^"]+)"/gi, '"<span class=sc_attr>$1</span>"' );
    /** Boolean/Null */
    string = string.replace( /\b(true|false|null)\b/gi, '<span class=sc_bool>$1</span>' );
    return string;
  }

  return module;
})({});

var editor = new EditorJS({
    
    readOnly: false,

    holder: 'story-editor',

    tools: {

      header: {
        class: Header,
        inlineToolbar: ['marker', 'link'],
        config: {
          placeholder: 'Header'
        },
        shortcut: 'CMD+SHIFT+H'
      },


      image: {
        class: ImageTool,
        config: {
          
          uploader: {

            uploadByFile(file){


              return mediaUploader(file).then(function(data){

                if(data.status != "success"){

                  return {
                    success: 0
                  }
                }

                return {
                  success: 1,
                  file: {

                    url: '/api/media/download/c/' + data.reply,
              
                  }
                }


              })

            }

          }

        }
      },

      list: {
        class: EditorjsList,
        inlineToolbar: true,
        shortcut: 'CMD+SHIFT+L'
      },

      checklist: {
        class: Checklist,
        inlineToolbar: true,
      },

      quote: {
        class: Quote,
        inlineToolbar: true,
        config: {
          quotePlaceholder: 'Enter a quote',
          captionPlaceholder: 'Quote\'s author',
        },
        shortcut: 'CMD+SHIFT+O'
      },

      warning: Warning,

      marker: {
        class:  Marker,
        shortcut: 'CMD+SHIFT+M'
      },

      code: {
        class:  CodeTool,
        shortcut: 'CMD+SHIFT+C'
      },

      delimiter: Delimiter,

      inlineCode: {
        class: InlineCode,
        shortcut: 'CMD+SHIFT+C'
      },

      linkTool: LinkTool,

      embed: Embed,

      table: {
        class: Table,
        inlineToolbar: true,
        shortcut: 'CMD+ALT+T'
      },

    },

    data: {
      blocks: [
        {
          type: "header",
          data: {
            text: "Title",
            level: 1
          }
        },
        {
          type: "header",
          data: {
            text: "[0000-00]",
            level: 3
          }
        },
        {
          type : 'paragraph',
          data : {
            text : 'Write your thing'
          }
        }
      ]
    },
    onReady: function(){
      console.log("data ready")
    },
//    onChange: function(api, event) {
//      console.log('something changed', event);
//    }
  }
);



const saveButton = document.getElementById('saveButton');


const toggleReadOnlyButton = document.getElementById('toggleReadOnlyButton');
const readOnlyIndicator = document.getElementById('readonly-state');


saveButton.addEventListener('click', async function () {
  
    let savedData = await editor.save()

    let a_info = JSON.parse(JSON.stringify(ARTICLE_INFO))

    a_info.title = savedData.blocks[0].data.text

    a_info.dateMarked = savedData.blocks[1].data.text

    blen = savedData.blocks.len 

    for (let i = 0 ; i < blen; i ++){




    }

    a_info.content = JSON.stringify(savedData)

    let req = {
        data: JSON.stringify(a_info)
    }

    let resp = await fetch("/api/story/upload",{
        body: JSON.stringify(req),
        headers: {
            'Content-Type': 'application/json'
            },
        method: "POST"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to submit content")

        return
    }

    alert("successfully submitted: " + result.reply)

});

toggleReadOnlyButton.addEventListener('click', async () => {
  const readOnlyState = await editor.readOnly.toggle();

  readOnlyIndicator.textContent = readOnlyState ? 'On' : 'Off';
});

