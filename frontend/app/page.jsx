'use client'
import React, { useState } from "react";
import { Tooltip } from 'react-tooltip';
import PropTypes from 'prop-types';
import Image from 'next/image';
function IconWithTooltip(props) {
    return (
        <div className="px-2 flex flex-grow-0">
            <a href={props.link}
                className="relative text-textsecondary hover:text-white group"
                data-tooltip-id={props.tooltipText + "-tooltip"}
            >
                {props.children}
            </a>
            <Tooltip
                id={props.tooltipText + "-tooltip"}
                content={props.tooltipText}
                opacity={1}
            />
        </div>
    );
}

IconWithTooltip.propTypes = {
    link: PropTypes.string.isRequired,
    tooltipText: PropTypes.string.isRequired,
    children: PropTypes.node.isRequired
}

function IntroductionCard() {
    const [emailTooltipContent, setEmailTooltipContent] = useState("");
    return (
        <div className="flex items-center">
            <div className="flex px-4">
                <Image src="/my-photo.webp" className="rounded-full h-100 pt-3" alt="my face" width="1024" height="1024"/>
            </div>
            <div className="py-4 w-full">
                <div className="text-5xl font-bold text-center">Hi, I am Vladimir Zagorovskii</div>
                <div className="text-4xl font-bold text-center">Back-End developer</div>
                <div className="text-2xl text-center text-textsecondary">18 y.o developer from Moscow</div>
                <div className="flex justify-center items-center pt-3">
                    <IconWithTooltip
                        tooltipText="Go to my LinkedIn"
                        link="https://www.linkedin.com/in/vladimir-zagorovskii-00b49a346/"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 504.4 504.4" fill="currentColor"><path d="M377.6,0.2H126.4C56.8,0.2,0,57,0,126.6v251.6c0,69.2,56.8,126,126.4,126H378c69.6,0,126.4-56.8,126.4-126.4V126.6C504,57,447.2,0.2,377.6,0.2z M168,408.2H96v-208h72V408.2z M131.6,168.2c-20.4,0-36.8-16.4-36.8-36.8c0-20.4,16.4-36.8,36.8-36.8c20.4,0,36.8,16.4,36.8,36.8C168,151.8,151.6,168.2,131.6,168.2z M408.4,408.2H408h-60V307.4c0-24.4-3.2-55.6-36.4-55.6c-34,0-39.6,26.4-39.6,54v102.4h-60v-208h56v28h1.6c8.8-16,29.2-28.4,61.2-28.4c66,0,77.6,38,77.6,94.4V408.2z"></path></svg>
                    </IconWithTooltip>
                    <IconWithTooltip
                        tooltipText="Go to my github"
                        link="https://github.com/svorogaze"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 16 16"><path fill="currentColor" d="M8 0c4.42 0 8 3.58 8 8a8.01 8.01 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38c0-.27.01-1.13.01-2.2c0-.75-.25-1.23-.54-1.48c1.78-.2 3.65-.88 3.65-3.95c0-.88-.31-1.59-.82-2.15c.08-.2.36-1.02-.08-2.12c0 0-.67-.22-2.2.82c-.64-.18-1.32-.27-2-.27s-1.36.09-2 .27c-1.53-1.03-2.2-.82-2.2-.82c-.44 1.1-.16 1.92-.08 2.12c-.51.56-.82 1.28-.82 2.15c0 3.06 1.86 3.75 3.64 3.95c-.23.2-.44.55-.51 1.07c-.46.21-1.61.55-2.33-.66c-.15-.24-.6-.83-1.23-.82c-.67.01-.27.38.01.53c.34.19.73.9.82 1.13c.16.45.68 1.31 2.69.94c0 .67.01 1.3.01 1.49c0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8" /></svg>
                    </IconWithTooltip>
                    <IconWithTooltip
                        tooltipText="Message me on Telegram"
                        link="https://t.me/svorogaze"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" viewBox="0 0 24 24"><path d="M9.862 15.897 9.826 20c.422 0 .607-.18.83-.396l1.987-1.91 4.118 3.01c.75.413 1.284.2 1.472-.697L21.93 4.47c.272-1.142-.417-1.593-1.157-1.325L2.57 10.35c-1.102.434-1.086 1.05-.198 1.333l4.438 1.383L17.3 6.615c.44-.299.837-.136.509.163" /></svg>
                    </IconWithTooltip>
                    <button className="relative text-textsecondary hover:text-text group px-2"
                        data-tooltip-id="email-tooltip"
                        onClick={() => {
                            navigator.clipboard.writeText("vvzagorovskiy@gmail.com");
                            setEmailTooltipContent("Email copied");
                        }}
                        onMouseEnter={() => {
                            setEmailTooltipContent("Copy email");
                        }}
                        onMouseLeave={() => {
                            setEmailTooltipContent("");
                        }}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path fill="currentColor" d="M4 20q-.825 0-1.412-.587T2 18V6q0-.825.588-1.412T4 4h16q.825 0 1.413.588T22 6v12q0 .825-.587 1.413T20 20zm8-7l8-5V6l-8 5l-8-5v2z" /></svg>
                    </button>
                    <Tooltip
                        id="email-tooltip"
                        content={emailTooltipContent}
                        opacity={1}
                    />
                </div>
            </div>
        </div>
    );
}

function SkillCard(props) {
    return (
        <div className="w-1/2 lg:w-1/5 md:w-1/3 xl:w-1/6 bg-background justify-center items-center flex">
            <div className="bg-foreground px-3 items-center">
                <Image
                    src={props.imageLink}
                    alt=""
                    className="object-cover h-48 w-48 p-3"
                    width="1024"
                    height="1024"
                />
                <div className="text-center text-2xl w-full h-full">{props.skillName}</div>
            </div>
        </div>
    );
}

function Skills() {
    const skills = [
        {
            skillName: "C++",
            imageLink: "/c++-icon.png"
        },
        {
            skillName: "Go",
            imageLink: "/golang-icon.png"
        },
        {
            skillName: "Python",
            imageLink: "/python-icon.png"
        },
        {
            skillName: "JS",
            imageLink: "/js-icon.png"
        },
        {
            skillName: "React",
            imageLink: "/react-icon.png"
        },
        {
            skillName: "SQL",
            imageLink: "/sql-icon.png"
        },
        {
            skillName: "Git",
            imageLink: "/git-icon.png"
        }
    ];
    return (
        <div className="w-full py-4">
            <div>
                <div className="py-4 w-full">
                    <div className="text-5xl font-bold text-center">Skills</div>
                </div>
            </div>
            <div className="flex flex-wrap justify-center items-stretch gap-6 py-2">
                {skills.map((e) => {
                    return <SkillCard skillName={e.skillName} imageLink={e.imageLink} />;
                })}
            </div>
        </div>
    );
}

function Achievement(props) {
    return (
        <div className="bg-foreground rounded-2xl shadow-lg p-6 max-w-xl w-full transition hover:shadow-xl">
            <h3 className="text-2xl font-bold mb-1">{props.text}</h3>
            <p className="text-text">{props.detail}</p>
        </div>
    );
}

function Achievements() {
    return (
        <div className="w-full py-10">
            <h2 className="text-5xl font-bold text-center mb-8">Achievements</h2>
            <div className="flex flex-col gap-6 items-center">
                <Achievement text="Silver Medal" detail="International Zhautykov Olympiad in Informatics 2024"/>
                <Achievement text="Top 70 nationally" detail="Finals of National Russian Olympiad in Informatics"/>
                <div className="bg-foreground rounded-2xl shadow-lg p-6 max-w-xl w-full transition hover:shadow-xl">
                    <h3 className="text-2xl font-bold mb-1">Codeforces</h3>
                    <p className="text-text">
                        Maximum rating of <span className="font-bold">1944</span> on{' '}
                        <a href="https://codeforces.com/profile/Lobotomy_Corporation_Fan"
                           className="text-blue-400 hover:text-blue-300 underline transition duration-150">
                            codeforces
                        </a>{' '}
                        (top <span className="font-bold">3%</span>)
                    </p>
                </div>
            </div>
        </div>
    );
}

export default function Home() {
    return (
        <div className="text-text">
            <IntroductionCard />
            <Skills />
            <Achievements/>
        </div>
    );
}