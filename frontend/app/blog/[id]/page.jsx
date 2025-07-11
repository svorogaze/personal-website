'use client'
import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Typography, Box, Avatar, Divider } from '@mui/material';
import { format } from 'date-fns';
import parse from 'html-react-parser';
function BlogHeader(props) {
    return (
    <Box sx={{ maxWidth: 1000, mx: 'auto', px: { xs: 2, md: 0 }, pt: 4}} className="text-text">
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

function BlogText(props) {

    return (
        <Box sx={{ maxWidth: 1000, mx: 'auto', px: { xs: 2, md: 0 }, pb: 4, lineHeight: 1.6}} style={{whiteSpace: 'pre-line'}} className="text-2xl text-text">
          {parse(props.text)}
        </Box>
    );
}


function Blog(props) {
    return (
        <div className="w-full items-center justify-center flex">
            <div className="h-full items-center justify-center px-4 text-text text-left">
                <BlogHeader title={props.title} author={props.author} date={props.date}/>
                <BlogText text={props.text}/>
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
        <div className="scrollbar-stable">
            <Blog title={blogData.data.title} author={blogData.metadata} date={blogData.data.creationDate} text={blogData.data.contents}/>
        </div>
    );
}