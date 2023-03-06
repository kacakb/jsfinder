<div align="center">
  <p>
    <img src="https://user-images.githubusercontent.com/64865400/223095605-38da9d6b-c9fa-4bfd-976a-8ed68a2812c2.png" alt="Logo">
  </p>
  <p>
    <a href="https://golang.org/doc/go1.20"><img src="https://img.shields.io/badge/Go-v1.20-blue"></a>
    <a href="https://github.com/kacakb/jsfinder/releases"><img src="https://img.shields.io/badge/releases-latest-brightgreen.svg"></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
    <a href="https://github.com/kacakb/jsfinder/issues"><img src="https://img.shields.io/badge/Issues-Welcome-blueviolet"></a>
  </p>
  <p>
    <a href="#features">Features</a> |
    <a href="#installation">Installation</a> |
    <a href="#usage">Usage</a> |
    <a href="#demo">Demo</a> |
    <a href="#contributing">Contributing</a> |
    <a href="#license">License</a> |
    <a href="#contact">Contact</a>
  </p>
</div>

jsFinder is a command-line tool written in Go that scans web pages to find JavaScript files linked in the HTML source code. It searches for any attribute that can contain a JavaScript file (e.g., src, href, data-main, etc.) and extracts the URLs of the files to a text file. The tool is designed to be simple to use, and it supports reading URLs from a file or from standard input.

JSFinder is useful for web developers and security professionals who want to find and analyze the JavaScript files used by a web application. By analyzing the JavaScript files, it's possible to understand the functionality of the application and detect any security vulnerabilities or sensitive information leakage.


<h2 id="features">Features</h2>

<ul>
  <li>Reading URLs from a file or from <strong>stdin</strong> using command line arguments.</li>
  <li>Running <strong>multiple</strong> HTTP GET requests concurrently to each URL.</li>
  <li>Limiting the <strong>concurrency</strong> of HTTP GET requests using a  flag.</li>
  <li>Using a <strong>regular expression</strong> to search for JavaScript files in the response body of the HTTP GET requests.</li>
  <li><strong>Writing the found JavaScript files</strong> to a file specified in the command line arguments or to a default file named "output.txt".</li>
  <li><strong>Printing informative messages to the console</strong> indicating the status of the program's execution and the output file's location.</li>
  <li>Allowing the program to run in <strong>verbose or silent mode</strong> using a flag.</li>
   </ul>
   
   <h2 id="installation">Installation</h2>
   
   jsfinder requires Go 1.20 to install successfully.Run the following command to get the repo :
   
  <pre><code class="language-go">go install -v github.com/kacakb/jsfinder@latest</code><button class="btn" data-clipboard-text="go install -v github.com/kacakb/jsfinder@latest"></button></pre>

<h2 id="usage">Usage</h2>
 <pre><code class="language-go">jsfinder -h </code><button class="btn" data-clipboard-text="go install -v github.com/kacakb/jsfinder@latest"></button></pre>
 
