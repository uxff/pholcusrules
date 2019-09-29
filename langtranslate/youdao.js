
var request = require("request");
var options = {
    url:'http://fanyi.youdao.com/translate?client=deskdict&keyfrom=chrome.extension&xmlVersion=1.1&dogVersion=1.0&ue=utf8&i=' + 'do%20you' + '&doctype=xml',
    headers: {
        'Cookie': 'OUTFOX_SEARCH_USER_ID=1646972453@111.199.190.14; _ntes_nnid=96dc4b0801e97f57fe83caa919bf51fd,1567345676386; OUTFOX_SEARCH_USER_ID_NCOO=1247166220.1406908'
    }
};

var parseXml =     (text, tagName) => {
    var parser = new window.DOMParser();
    var xml = parser.parseFromString(text.replace(/&lt;/g, '<').replace(/&gt;/g, '>'), 'text/xml');
    var value = {};

    tagName = Array.isArray(tagName) ? tagName : [tagName];

    tagName.map(e => {
        value[e] = xml.getElementsByTagName(e)[0].childNodes[0].nodeValue;
    });

    return { xml, value };
};



request.get(options, function(err, response, body){
    console.info(response.body);
    // var xml2js = require('xml2js');
    // var xmlParser = new xml2js.Parser({explicitArray: false, ignoreAttrs: true});
    // xmlParser.parseString(response, function(err, result) {
    //     //message.ContentTrans = 'transed+' + result.translation;
    //     console.log('err when parse=' + err + ' translation=%o' + result);
    // });
    let res = parseXml(response.body, 'translation');
    // var domParser=require('xmldom').DOMParser;
    // var parser=new domParser();
    // var parent = parser.parseFromString(response.body);

    // var tempVal=parent.getElementsByTagName("translation")[0].nodeValue;//[0].childNodes[0].nodeValue;
    console.log("tempVal:"+res);

});







