"use client";

import * as React from "react";
import {
    ArrowUpCircleIcon,
    BarChartIcon,
    ClipboardListIcon,
    DatabaseIcon,
    FileIcon,
    FolderIcon,
    HelpCircleIcon,
    LayoutDashboardIcon,
    SearchIcon,
    SettingsIcon,
    UsersIcon,
} from "lucide-react";

import { NavDocuments } from "@/components/nav-documents";
import { NavMain } from "@/components/nav-main";
import { NavSecondary } from "@/components/nav-secondary";
import { NavUser } from "@/components/nav-user";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link } from "react-router";

const data = {
    user: {
        name: "shadcn",
        email: "m@example.com",
        avatar: "https://png.pngtree.com/png-vector/20220817/ourmid/pngtree-women-cartoon-avatar-in-flat-style-png-image_6110776.png",
    },
    navMain: [
        {
            title: "Dashboard",
            url: "/",
            icon: LayoutDashboardIcon,
        },
        {
            title: "Analytics",
            url: "/analytics",
            icon: BarChartIcon,
        },
        {
            title: "Projects",
            url: "/projects",
            icon: FolderIcon,
        },
        {
            title: "Customers",
            url: "/customers",
            icon: UsersIcon,
        },
    ],
    navSecondary: [
        {
            title: "Settings",
            url: "/settings",
            icon: SettingsIcon,
        },
        {
            title: "Get Help",
            url: "/help",
            icon: HelpCircleIcon,
        },
        {
            title: "Search",
            url: "/search",
            icon: SearchIcon,
        },
    ],
    documents: [
        {
            name: "Data Library",
            url: "/data-library",
            icon: DatabaseIcon,
        },
        {
            name: "Reports",
            url: "/reports",
            icon: ClipboardListIcon,
        },
        {
            name: "Word Assistant",
            url: "/word-assistant",
            icon: FileIcon,
        },
    ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    return (
        <Sidebar collapsible="offcanvas" {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton
                            asChild
                            className="data-[slot=sidebar-menu-button]:!p-1.5"
                        >
                            <Link to="/">
                                <ArrowUpCircleIcon className="h-5 w-5" />
                                <span className="text-base font-semibold">
                                    Acme Inc.
                                </span>
                            </Link>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <NavMain items={data.navMain} />
                <NavDocuments items={data.documents} />
                <NavSecondary items={data.navSecondary} className="mt-auto" />
            </SidebarContent>
            <SidebarFooter>
                <NavUser user={data.user} />
            </SidebarFooter>
        </Sidebar>
    );
}
