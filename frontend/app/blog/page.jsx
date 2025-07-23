'use client'
import React, { useState, use, useEffect } from "react";
import { useSearchParams, usePathname, useRouter} from 'next/navigation';
import { useDebouncedCallback } from 'use-debounce';
import {Pagination} from "@heroui/react";
import Image from 'next/image';
import Link from 'next/link';
function BlogCard(props) {
    return (
        <div className="w-1/2 md:w-1/3 lg:w-1/4 xl:w-1/5 px-3 py-3">
            <Link className="blog-card h-full w-full flex flex-col max-w-full group" href={`/blog/${props.id}`}>
                <Image
                    src={`${window.location.origin}/images/blog-cover-images/` + props.imageLink}
                    alt=""
                    width="512"
                    height="512"
                    className="object-cover w-full h-48 group-hover:brightness-50 transition-all duration-75"
                />
                <div className="px-5 py-4 text-left flex-grow max-w-full">
                    <div className="font-bold text-xl mb-2 pr-2 max-w-full break-words">{props.title}</div>
                    <p className="text-base break-words">{props.text}</p>
                </div>
            </Link>
        </div>
    );
}

function CardList(props) {
    return (
        <div className="flex flex-wrap justify-center items-stretch gap-1">
            {props.blogs.map((e) => {
                return <BlogCard key={e.id} title={e.title} text={e.description} id={e.id} imageLink={e.pictureLink} />;
            })}
        </div>
    );
}


function SearchBar(props) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();
    const HandleInput = useDebouncedCallback((name) => {
        const params = new URLSearchParams(searchParams);
        params.set('page', '1');
        if (name) {
            params.set("query", name);
        } else {
            params.delete("query");
        }
        replace(`${pathname}?${params.toString()}`);
    }, 100);
    return (
        <div className="w-full pt-4 pb-1 px-1">
            <div className="px-1">
                <div className="flex justify-center items-center">
                    <input type="text" className="h-14 w-96 pr-8 pl-5 rounded z-0 focus:outline-none bg-foreground text-text placeholder-textsecondary"
                        placeholder={props.placeholder}
                        value={props.searchText}
                        onChange={(e) => HandleInput(e.target.value)}
                        defaultValue={searchParams.get('query')?.toString()}
                    >
                    </input>
                </div>
            </div>
        </div>
    );
}

function BlogPagination(props) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();
    const HandlePage = useDebouncedCallback((number) => {
        const params = new URLSearchParams(searchParams);
        params.set("page", number);
        replace(`${pathname}?${params.toString()}`);
    }, 200);
    return (
        <div className="flex items-center justify-center pt-3 pb-1">
            <Pagination initialPage={props.currentPage} total={props.allPages} showControls onChange={(n) => HandlePage(n)}/>
        </div>
    );
}

export default function Home(props) {
    const params = use(props.searchParams);
    const query = params.query || "";
    const currentPage = Number(params.page) || 1;
    const [blogData, setBlogData] = useState();
    useEffect(() => {
        const fetchData = async () => {
            const url = `${window.location.origin}/api/blogs?offset=${(currentPage - 1) * 8}&size=${8}&search-query=${query}`;
            const response = await fetch(url);
            const data = await response.json();
            setBlogData(data);
        };
        fetchData();
    }, [query, currentPage]);
    if (blogData == undefined) {
        return (
            <div className="scrollbar-stable"/>
        );
    }
    const allPages = (blogData.metadata.size + 7) / 8;
    console.log(blogData);
    return (
        <div className="scrollbar-stable">
            <SearchBar placeholder="Look up anything..."/>
            <BlogPagination currentPage={currentPage} allPages={allPages}/>
            <CardList blogs={blogData.data}/>
        </div>
    );
}
