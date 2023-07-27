package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="...">) in an HTML
// document.
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a
// slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	//html.parse returns nodes of the html file
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	//we are passing the html nodes as a pointer here because 
	//efficiently modify and traverse
	// the tree-like structure representing the HTML document.
	nodes := linkNodes(doc)
	//linkNodes returns a slice pointers of link nodes which data is a <a>
	var links []Link
	for _, node := range nodes {
		//we are extracting the href and text inside each link in the
		//html file.
		//we are traversing the dom treeto extract each href and text content
		// for each link found
		links = append(links, buildLink(node))
	}
	// returns the href and text struct
	return links, nil
}

//for extracting the href in the link and text for a node link
func buildLink(n *html.Node) Link {
	var ret Link
	//n.Attr, = Attr is  an array of struct field in the Node struct, attr represents the attr
	//a node has e.g href for <a> tag and class names and so on, so we are 
	//extracting it
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

//this is the function to extract the texts from the node
func text(n *html.Node) string {
	//Node type is an indicator used to identify the type of node its not part
	//of the node struct
	/*
	const (
    ErrorNode NodeType = iota // Node is an error node
    TextNode                  // Node is a text node
    DocumentNode              // Node is the document itself
    ElementNode               // Node is an HTML element
    CommentNode               // Node is a comment
    DoctypeNode               // Node is a DOCTYPE declaration
)
	*/
	if n.Type == html.TextNode {
		//type and data is a part of the node struct, 
		//we are checking here if the type of the node is a text,
		//then returning the data =<a>
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	//here is where we run the recursive loop to traverse the whole dom tree for a node
	//so since n would be a link tag <a> we are setting c=firstchild
	//which is a text and while c is not empty we are increasing it to the next sibling
	//which is another link tag<a>
	//first child of the link is a text
	//then we are storing the text in ret 
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}


//here we are collecting each node which has a <a> tag 
func linkNodes(n *html.Node) []*html.Node {
	//checming if the type of the node is an element
	// and the data of the node is <a>
	//then we arereturning the slice of node struct with a pointer consisting of
	//the node slice would be an slice of nodes that is an <a> element
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	//traversing the html document for each <a> node tag
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}