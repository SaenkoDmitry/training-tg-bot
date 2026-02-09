import DOMPurify from 'dompurify';

export default function SafeTextRenderer({html}) {
    const cleanHTML = DOMPurify.sanitize(html);
    return <div dangerouslySetInnerHTML={{__html: cleanHTML}}/>;
}
