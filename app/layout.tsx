import { AppSidebar } from "~/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "~/components/ui/sidebar";
import { SiteHeader } from "~/components/site-header";
import { Outlet } from "react-router";

export default function Layout() {
    return (
        <SidebarProvider>
            <AppSidebar variant="inset" />
            <SidebarInset>
                <SiteHeader />
                <div className="flex flex-1 flex-col">
                    <div className="@container/main flex flex-1 flex-col gap-2">
                        <Outlet />
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}
