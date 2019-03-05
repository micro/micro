const randomElement = (arr = []) => {
    return arr[Math.floor(Math.random() * arr.length)];
};

const kebab = (str: string) => {
    return (str || "").replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
};

const toggleFullScreen = () => {
    let doc = window.document;
    let docEl = doc.documentElement;

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

const setCookie = (a: string, v: string, days: number) => {
    let expires = "";
    if (days) {
        let date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = a + "=" + (v || "") + expires + "; path=/";
}

export default {
    randomElement,
    toggleFullScreen,
    kebab,
    getCookieValue,
    setCookie
};
