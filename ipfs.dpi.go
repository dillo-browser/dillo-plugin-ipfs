package main

import (
    "bufio"
    "fmt"
    "html"
    "io"
    "net"
    "net/http"
    "net/url"
    "os"
    "os/user"
    "path"
    "strings"
)

var EndOfTag = fmt.Errorf("end of tag")

type Tag struct {
    Props map[string]string
}

func main() {
    fmt.Printf("[ipfs dpi]: Starting\n")
    listener, err := net.FileListener(os.Stdin)
    if err != nil {
        fmt.Fprintf(os.Stderr, "listener: %s\n", err.Error())
        os.Exit(1)
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Fprintf(os.Stderr, "accept: %s\n", err.Error())
            continue
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

    for {
        err := handleTag(conn, reader)
        if (err != nil) {
            if (err == io.EOF) {
                return
            }
            fmt.Fprintf(os.Stderr, "[ipfs dpi]: read error: %s\n", err.Error())
            return
        }
    }
}

func handleTag(conn net.Conn, reader *bufio.Reader) error {
    tag, err := readTag(reader)
    if (err != nil) {
        return err
    }
    cmd := tag.Props["cmd"]
    switch cmd {
    case "auth": return handleAuth(conn, tag.Props["msg"])
    case "open_url": return handleOpenUrl(conn, tag.Props["url"])
    default: return fmt.Errorf("unhandled cmd \"%s\"\n", cmd)
    }
}

func handleAuth(conn net.Conn, msg string) error {
    usr, err := user.Current()
    if err != nil {
        return err
    }

    filename := path.Join(usr.HomeDir, ".dillo", "dpid_comm_keys")
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    var pid int
    var key string
    _, err = fmt.Fscanf(file, "%d %s", &pid, &key)
    if err != nil {
        return err
    }
    _ = pid
    if (msg != key) {
        return fmt.Errorf("mismatched dpid key %s/%s", msg, key)
    }
    return nil
}

func handleOpenUrl(conn net.Conn, url_str string) error {
    u, err := url.Parse(url_str)
    if err != nil {
        return err
    }
    switch u.Scheme {
    case "ipfs", "ipns": return serveIpfs(conn, url_str, u)
    default: return serve404(conn, url_str)
    }
}

func escapeDpiValue(str string) string {
    return strings.Replace(str, "'", "''", -1)
}

func writeHeader(conn net.Conn, url_str string, mime_type string) error {
	_, err := fmt.Fprintf(conn, "<cmd='start_send_page' url='%s' '>Content-Type: %s\r\n\r\n", escapeDpiValue(url_str), mime_type)
	return err
}

func writeStatus(conn net.Conn, msg string) error {
	_, err := fmt.Fprintf(conn, "<cmd='send_status_message' msg='%s' '>", escapeDpiValue(msg))
	return err
}

func serve404(conn net.Conn, url_str string) error {
    err := writeHeader(conn, url_str, "text/html")
    if err != nil {
        return err
    }
    _, err = fmt.Fprintf(conn, "<h3>Not Found</h3>")
    if err != nil {
        return err
    }
    return io.EOF
}

func IpfsUrlToGatewayUrl(u *url.URL) string {
    var prefix string
    if (u.Host == "") {
        prefix = ""
    } else {
        prefix = "/" + u.Host
    }
    u.Path = u.Scheme + prefix + u.Path
    u.Scheme = "http"
    u.Host = "127.0.0.1:8080"
    return u.String()
}

func serveIpfs(conn net.Conn, url_str string, u *url.URL) error {
    gwUrl := IpfsUrlToGatewayUrl(u)
    err := writeStatus(conn, "Fetching IPFS content...")
    if err != nil {
        return err
    }
    resp, err := http.Get(gwUrl)
    if err != nil {
        err := writeHeader(conn, url_str, "text/html")
        _, err = fmt.Fprintf(conn, "<h3>Gateway Error</h3>\n<pre>%s</pre>", html.EscapeString(err.Error()))
        if err != nil {
            return err
        }
        conn.Close()
        return io.EOF
    }
    defer resp.Body.Close()
    contentType := resp.Header.Get("Content-Type")
    if (contentType == "") {
        contentType = "text/plain"
    }
    err = writeHeader(conn, url_str, contentType)
    if err != nil {
        return err
    }
    var buf [512]byte
    for {
        n, err1 := resp.Body.Read(buf[0:])
        if err1 != nil && err1 != io.EOF {
            return err1
        }
        if (n > 0) {
            _, err2 := conn.Write(buf[0:n])
            if err2 != nil {
                return err
            }
        }
        if (err1 == io.EOF) {
            break
        }
    }
    return io.EOF
}

func readTag(reader *bufio.Reader) (*Tag, error) {
    c, err := reader.ReadByte()
    if err != nil {
        return nil, err
    }
    if c != '<' {
        return nil, fmt.Errorf("expected '<' but got '%c'\n", c)
    }
    props := make(map[string]string)
    for {
        key, value, err := readProperty(reader)
        if (err != nil) {
            if (err == EndOfTag) {
                break
            }
            return nil, err
        }
        props[key] = value
    }
    return &Tag{Props: props}, nil
}

func readProperty(reader *bufio.Reader) (string, string, error) {
    c1, err := reader.ReadByte()
    if err != nil {
        return "", "", err
    }
    if c1 == '\'' {
        c, err := reader.ReadByte()
        if err != nil {
            return "", "", err
        }
        if c == '>' {
            return "", "", EndOfTag
        } else {
            return "", "", fmt.Errorf("expected '>'")
        }
    }

    key_bytes := []byte{c1}
    key_bytes1, err := reader.ReadBytes('=')
    if err != nil {
        return "", "", err
    }
    key_bytes = append(key_bytes, key_bytes1[0:len(key_bytes1)-1]...)
    c, err := reader.ReadByte()
    if err != nil {
        return "", "", err
    }
    if c != '\'' {
        return "", "", fmt.Errorf("expected quote")
    }
    var value_bytes []byte
    for {
        part, err := reader.ReadSlice('\'')
        if err != nil {
            return "", "", err
        }
        value_bytes = append(value_bytes, part[0:len(part)-1]...)
        c, err := reader.ReadByte()
        if err != nil {
            return "", "", err
        }
        if c == '\'' {
            // escaped quote
            value_bytes = append(value_bytes, '\'')
        } else if c == ' ' {
            // end of value
            break
        } else {
            return "", "", fmt.Errorf("expected space but got '%c'", c)
        }
    }

    return string(key_bytes), string(value_bytes), nil
}

// <cmd='auth' msg='a20a4710' '>
// <cmd='open_url' url='ipfs://asdasdas/''asdasd/asdad>aasd' '>

