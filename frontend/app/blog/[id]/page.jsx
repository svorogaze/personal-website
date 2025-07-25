'use client'
import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Typography, Box, Avatar, Divider } from '@mui/material';
import { format } from 'date-fns';
import {ReactMarkdown, MarkdownHooks} from "react-markdown";
import remarkGfm from 'remark-gfm';;
import rehypeRaw from 'rehype-raw';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { darcula } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { visit } from 'unist-util-visit';
import rehypeStarryNight from 'rehype-starry-night';
function BlogHeader(props) {
    return (
    <Box sx={{ mx: 'auto', px: { xs: 2, md: 0 }, pt: 4}} className="text-text max-w-[min(1000px,100%)]">
      <Typography 
        variant="h3" 
        component="h1" 
        sx={{ 
          fontWeight: 700, 
          mb: 3,
          fontSize: { xs: '2rem', md: '2.5rem' },
          lineHeight: 1.2
        }}
      >
        {props.title}
      </Typography>
      
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
        <Avatar 
          alt={props.author.name} 
          src={props.author.avatarLink} 
          sx={{ width: 64, height: 64, mr: 2 }}
        />
        <Box>
          <Typography variant="subtitle1" component="div" sx={{ fontWeight: 600 }}>
            {props.author.name}
          </Typography>
          <Typography variant="body2" className="text-textsecondary">
            {format(props.date, 'MMMM d, yyyy')}
          </Typography>
        </Box>
      </Box>
      
      <Divider sx={{ mb: 4 }} />
    </Box>
  );
}

function MarkdownRenderer({ content }){
  return (
    <div className="markdown-body font-sans text-text max-w-[min(1000px,100%)] mx-auto py-8 text-2xl [&>*]:max-w-full">
      <MarkdownHooks
        remarkPlugins={
          [remarkGfm,
            () => (tree) => { // fix inline code
                visit(tree, 'code', (node) => {
                node.lang = node.lang ?? 'plaintext';
              });
            },
          ]
        }
        rehypePlugins={[rehypeRaw, rehypeStarryNight]}
        components={{
          h1: ({ node, ...props }) => (
            <h1 className="text-4xl lg:text-6xl font-bold border-b border-gray-200 pb-2 mt-6 mb-4" {...props} />
          ),
          h2: ({ node, ...props }) => (
            <h2 className="text-3xl lg:text-5xl font-bold border-b border-gray-200 pb-2 mt-5 mb-3" {...props} />
          ),
          h3: ({ node, ...props }) => (
            <h3 className="text-2xl lg:text-4xl font-bold mt-4 mb-2" {...props} />
          ),
          h4: ({ node, ...props }) => (
            <h4 className="text-xl lg:text-3xl font-bold mt-3 mb-2" {...props} />
          ),
          h5: ({ node, ...props }) => (
            <h5 className="text-lg lg:text-2xl font-bold mt-3 mb-2" {...props} />
          ),
          h6: ({ node, ...props }) => (
            <h5 className="text-base lg:text-xl font-bold mt-3 mb-2" {...props} />
          ),
          p: ({ node, ...props }) => (
            <p className="mb-4 leading-relaxed text-base lg:text-xl" {...props} />
          ),
          a: ({ node, ...props }) => (
            <a className="text-[#4493d4] underline text-base lg:text-xl [&>*]:text-inherit" {...props} />
          ),
          ul: ({ node, ...props }) => (
            <ul className="list-disc pl-6 mb-4 text-base lg:text-xl" {...props} />
          ),
          ol: ({ node, ...props }) => (
            <ol className="list-decimal pl-6 mb-4 text-base lg:text-xl" {...props} />
          ),
          li: ({ node, ...props }) => (
            <li className="mb-2 text-base lg:text-xl" {...props} />
          ),
          blockquote: ({ node, ...props }) => (
            <blockquote className="border-l-4 border-gray-300 pl-4 italic text-text mb-4 text-base lg:text-xl" {...props} />
          ),
          code: ({ node, inline, className, children, ...props }) => {
            const match = /language-(\w+)/.exec(className || '');
            return match ? (
              <div className="rounded-md overflow-hidden mb-4 text-sm lg:text-lg max-w-full">
                <SyntaxHighlighter
                  language={match?.[1] || 'text'}
                  style={darcula}
                  codeTagProps={{
                    style: {
                      fontFamily: 'ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace'
                    }
                  }}
                  {...props}
                >
                  {String(children).replace(/\n$/, '')}
                </SyntaxHighlighter>
              </div>
            ) : (
              <code className="bg-[#2b2b2b] rounded px-1 py-0.5 text-base lg:text-xl font-mono text-[#e8e8e8]" {...props}>
                {children}
              </code>
            );
          },
          pre: ({ node, ...props }) => (
            <div className="mb-4 rounded-md overflow-hidden text-base lg:text-xl" {...props} />
          ),
          table: ({ node, ...props }) => (
            <table className="border-collapse border border-gray-300 mb-4 w-full text-base lg:text-xl" {...props} />
          ),
          thead: ({ node, ...props }) => (
            <thead className="bg-foreground text-base lg:text-xl" {...props} />
          ),
          th: ({ node, isHeader, ...props }) => (
            <th className="border border-foreground px-3 py-2 font-semibold text-base lg:text-xl" {...props} />
          ),
          td: ({ node, isHeader, ...props }) => (
            <td className="border border-foreground px-3 py-2 text-base lg:text-xl"  {...props} />
          ),
          img: ({ node, ...props }) => (
            <img className="h-auto rounded-md border border-foreground my-4 max-w-full" loading="lazy" {...props} />
          ),
          hr: ({ node, ...props }) => (
            <hr className="border-t border-foreground my-6 text-base lg:text-xl" {...props} />
          ),
        }}
      >
        {content}
      </MarkdownHooks>
    </div>
  );
};

function Blog(props) {
    return (
        <div className="w-full items-center justify-center flex">
            <div className="max-w-[min(1000px,100%)] h-full items-center justify-center text-text text-left">
                <BlogHeader title={props.title} author={props.author} date={props.date}/>
                <MarkdownRenderer content={props.text}/>
            </div>
        </div>
    )
}


export default function Home() {
    const params = useParams();
    const { id } = params;
    const [blogData, setBlogData] = useState();
    useEffect(() => {
        const fetchData = async () => {
            const url = `${window.location.origin}/api/blogs/${id}`;
            const response = await fetch(url);
            const data = await response.json();
            console.log(data);
            setBlogData(data);
        };
        fetchData();
    }, [id]);
    if (blogData == undefined) {
        return (
            <div className="scrollbar-stable"/>
        );
    }
    console.log(blogData);
    return (
      <>
        <div className="scrollbar-stable px-4">
            <Blog title={blogData.data.title} author={blogData.metadata} date={blogData.data.creationDate} text={blogData.data.contents}/>
        </div>
      </>
    );
}