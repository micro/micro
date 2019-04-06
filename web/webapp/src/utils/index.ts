const toggleFullScreen = (divId?: string) => {


    let docEl = null;
    let doc = window.document;

    if (!divId) {
        docEl = doc.documentElement;
    } else {
        docEl = document.getElementById(divId)
    }


    // @ts-ignore
    let requestFullScreen = docEl.requestFullscreen || docEl.mozRequestFullScreen || docEl.webkitRequestFullScreen || docEl.msRequestFullscreen;

    // @ts-ignore
    let cancelFullScreen = doc.exitFullscreen || doc.mozCancelFullScreen || doc.webkitExitFullscreen || doc.msExitFullscreen;


    // @ts-ignore
    if (!doc.fullscreenElement && !doc.mozFullScreenElement && !doc.webkitFullscreenElement && !doc.msFullscreenElement) {
        requestFullScreen.call(docEl);
    } else {
        cancelFullScreen.call(doc);
    }
};

const getCookieValue = (a: string) => {
    let b = document.cookie.match('(^|;)\\s*' + a + '\\s*=\\s*([^;]+)');
    return b ? b.pop() : '';
}

const getDefaultLan = () => {
    let locale = getCookieValue('locale')
    return locale ? locale : 'en';
}

const setCookie = (a: string, v: string, days: number) => {
    let expires = "";
    if (days) {
        let date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = a + "=" + (v || "") + expires + "; path=/";
}


const copyTxt = (text: string, callback: Function) => {

    // @ts-ignore
    if (!navigator.clipboard) {
        return callback(fallbackCopyTextToClipboard(text))
    }

    let flag = false;

    // @ts-ignore
    navigator.clipboard.writeText(text).then(
        function () {
            callback(true)
        },
        function () {
            callback(false)
        },
    );
}


const secondsToHHMMSS = (secondsInput: number) => {
    let hours = Math.floor(secondsInput / 3600);
    secondsInput %= 3600;
    let minutes = Math.floor(secondsInput / 60);
    let seconds = secondsInput % 60;

    return hours + ':' + (minutes < 10 ? '0' + minutes : minutes) + ':' + (seconds < 10 ? '0' + seconds : seconds)
}

let _xools = {
    toggleFullScreen,
    getCookieValue,
    getDefaultLan,
    secondsToHHMMSS,
    setCookie,
    copyTxt,
}

let XTools = {}

// @ts-ignore
XTools.install = function (Vue, options) {

    // add the instance method
    if (!Vue.prototype.hasOwnProperty('$xools')) {
        Object.defineProperty(Vue.prototype, '$xools', {
            get: function get() {
                return _xools;
            },
        });
    }
}

export default {
    XTools, Utils: _xools
};


function fallbackCopyTextToClipboard(text: string) {
    let textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    let flag = false;
    try {
        document.execCommand('copy');
        flag = true
    } catch (err) {
        flag = false
    }

    document.body.removeChild(textArea);
    return flag;
}
