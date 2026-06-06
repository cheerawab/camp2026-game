"use client"

import {
  BellIcon,
  BookOpenIcon,
  ChevronDownIcon,
  ChevronsUpDownIcon,
  CircleHelpIcon,
  ClipboardListIcon,
  CopyIcon,
  DownloadIcon,
  FilterIcon,
  HomeIcon,
  InboxIcon,
  LayersIcon,
  LayoutDashboardIcon,
  LineChartIcon,
  MailIcon,
  MoreHorizontalIcon,
  PlusIcon,
  RadioTowerIcon,
  SaveIcon,
  SearchIcon,
  SettingsIcon,
  ShieldAlertIcon,
  SparklesIcon,
  StarIcon,
  Trash2Icon,
  UserIcon,
} from "lucide-react"
import type { ReactNode } from "react"
import { useForm } from "react-hook-form"
import {
  Area,
  AreaChart as RechartsAreaChart,
  CartesianGrid,
  XAxis,
} from "recharts"
import { toast } from "sonner"

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/accordion"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogMedia,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/shared/ui/alert-dialog"
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert"
import { AspectRatio } from "@/shared/ui/aspect-ratio"
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/avatar"
import { Badge } from "@/shared/ui/badge"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/shared/ui/breadcrumb"
import {
  ButtonGroup,
  ButtonGroupSeparator,
  ButtonGroupText,
} from "@/shared/ui/button-group"
import { Button } from "@/shared/ui/button"
import { Calendar } from "@/shared/ui/calendar"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card"
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/shared/ui/carousel"
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/ui/chart"
import { Checkbox } from "@/shared/ui/checkbox"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/shared/ui/collapsible"
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/shared/ui/combobox"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
} from "@/shared/ui/command"
import {
  ContextMenu,
  ContextMenuCheckboxItem,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuLabel,
  ContextMenuSeparator,
  ContextMenuShortcut,
  ContextMenuTrigger,
} from "@/shared/ui/context-menu"
import { DirectionProvider } from "@/shared/ui/direction"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/ui/dialog"
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/shared/ui/drawer"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/shared/ui/dropdown-menu"
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/shared/ui/empty"
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSeparator,
  FieldSet,
  FieldTitle,
} from "@/shared/ui/field"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/ui/form"
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/shared/ui/hover-card"
import { IconBadge } from "@/shared/ui/icon-badge"
import { InfoTile } from "@/shared/ui/info-tile"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
  InputGroupText,
  InputGroupTextarea,
} from "@/shared/ui/input-group"
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/shared/ui/input-otp"
import { Input } from "@/shared/ui/input"
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemFooter,
  ItemGroup,
  ItemHeader,
  ItemMedia,
  ItemSeparator,
  ItemTitle,
} from "@/shared/ui/item"
import { Kbd, KbdGroup } from "@/shared/ui/kbd"
import { Label } from "@/shared/ui/label"
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarSeparator,
  MenubarShortcut,
  MenubarTrigger,
} from "@/shared/ui/menubar"
import { MetricCard } from "@/shared/ui/metric-card"
import {
  NativeSelect,
  NativeSelectOptGroup,
  NativeSelectOption,
} from "@/shared/ui/native-select"
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/shared/ui/navigation-menu"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/shared/ui/pagination"
import {
  Popover,
  PopoverContent,
  PopoverDescription,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/shared/ui/popover"
import { Progress } from "@/shared/ui/progress"
import { RadioGroup, RadioGroupItem } from "@/shared/ui/radio-group"
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/shared/ui/resizable"
import { ScrollArea, ScrollBar } from "@/shared/ui/scroll-area"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/select"
import { Separator } from "@/shared/ui/separator"
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/shared/ui/sheet"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSkeleton,
  SidebarProvider,
  SidebarTrigger,
} from "@/shared/ui/sidebar"
import { Skeleton } from "@/shared/ui/skeleton"
import { Slider } from "@/shared/ui/slider"
import { Toaster } from "@/shared/ui/sonner"
import { Spinner } from "@/shared/ui/spinner"
import { StatusBadge } from "@/shared/ui/status-badge"
import { StatusDot } from "@/shared/ui/status-dot"
import { Switch } from "@/shared/ui/switch"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/ui/tabs"
import { Textarea } from "@/shared/ui/textarea"
import { ToggleGroup, ToggleGroupItem } from "@/shared/ui/toggle-group"
import { Toggle } from "@/shared/ui/toggle"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/shared/ui/tooltip"

type DemoFormValues = {
  codename: string
}

const chartData = [
  { day: "Day 1", score: 18 },
  { day: "Day 2", score: 32 },
  { day: "Day 3", score: 28 },
  { day: "Day 4", score: 46 },
  { day: "Day 5", score: 54 },
]

const chartConfig = {
  score: {
    label: "Score",
    color: "var(--primary)",
  },
} satisfies ChartConfig

const sections = [
  "Foundations",
  "Forms",
  "Overlays",
  "Data",
  "Navigation",
  "Layout",
  "Project UI",
] as const

export function ComponentsPage() {
  return (
    <TooltipProvider>
      <main className="bg-background text-foreground min-h-svh">
        <Toaster position="top-right" />
        <div className="mx-auto flex w-full max-w-[1120px] flex-col gap-8 px-4 py-5 sm:px-6 lg:py-8">
          <header className="bg-surface-raised border-ink flex flex-col gap-5 rounded-[2rem] border-2 p-5 shadow-[6px_6px_0_rgba(23,35,58,0.16)] sm:p-7 lg:flex-row lg:items-end lg:justify-between">
            <div className="space-y-3">
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem>
                    <BreadcrumbLink href="/">Home</BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator />
                  <BreadcrumbItem>
                    <BreadcrumbPage>Components</BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
              <div className="space-y-2">
                <Badge
                  variant="secondary"
                  className="border-ink w-fit border-2 font-black"
                >
                  FIELD KIT UI
                </Badge>
                <h1 className="text-3xl font-black tracking-normal sm:text-5xl">
                  Camp 2026 元件規格板
                </h1>
                <p className="text-muted-foreground max-w-2xl text-sm sm:text-base">
                  對齊 Field Kit Collectible 視覺方向：紙感底色、粗邊框、
                  收藏卡厚影與現場工具感。
                </p>
              </div>
            </div>
            <nav className="flex flex-wrap gap-2">
              {sections.map((section) => (
                <Button
                  key={section}
                  variant="outline"
                  size="sm"
                  className="border-ink bg-card border-2 font-black shadow-[2px_2px_0_rgba(23,35,58,0.16)]"
                  asChild
                >
                  <a href={`#${section.toLowerCase().replace(" ", "-")}`}>
                    {section}
                  </a>
                </Button>
              ))}
            </nav>
          </header>

          <GallerySection
            id="foundations"
            title="Foundations"
            description="Buttons, feedback, surfaces, loading states, and small primitives."
          >
            <Preview title="Accordion">
              <Accordion type="single" collapsible className="w-full">
                <AccordionItem value="mission">
                  <AccordionTrigger>Mission briefing</AccordionTrigger>
                  <AccordionContent>
                    Keep the base UI small, composable, and easy to scan.
                  </AccordionContent>
                </AccordionItem>
                <AccordionItem value="reward">
                  <AccordionTrigger>Reward table</AccordionTrigger>
                  <AccordionContent>
                    Badges and status chips map to semantic color tokens.
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </Preview>

            <Preview title="Alert">
              <Alert>
                <BellIcon />
                <AlertTitle>System ready</AlertTitle>
                <AlertDescription>
                  Health check and player session panels can use the same
                  feedback surface.
                </AlertDescription>
              </Alert>
            </Preview>

            <Preview title="AspectRatio">
              <AspectRatio
                ratio={16 / 9}
                className="bg-muted flex items-center justify-center rounded-md border"
              >
                <SparklesIcon className="text-muted-foreground size-10" />
              </AspectRatio>
            </Preview>

            <Preview title="Avatar">
              <div className="flex items-center gap-3">
                <Avatar>
                  <AvatarImage src="https://github.com/shadcn.png" alt="" />
                  <AvatarFallback>SC</AvatarFallback>
                </Avatar>
                <Avatar>
                  <AvatarFallback>26</AvatarFallback>
                </Avatar>
                <Avatar>
                  <AvatarFallback>
                    <UserIcon className="size-4" />
                  </AvatarFallback>
                </Avatar>
              </div>
            </Preview>

            <Preview title="Badge">
              <div className="flex flex-wrap gap-2">
                <Badge>Default</Badge>
                <Badge variant="secondary">Secondary</Badge>
                <Badge variant="outline">Outline</Badge>
                <Badge variant="destructive">Destructive</Badge>
              </div>
            </Preview>

            <Preview title="Button">
              <div className="flex flex-wrap gap-2">
                <Button>
                  <SaveIcon />
                  Save
                </Button>
                <Button variant="secondary">
                  <CopyIcon />
                  Copy
                </Button>
                <Button variant="outline">
                  <DownloadIcon />
                  Export
                </Button>
                <Button variant="ghost" size="icon" aria-label="More">
                  <MoreHorizontalIcon />
                </Button>
              </div>
            </Preview>

            <Preview title="Card">
              <Card className="max-w-sm shadow-none">
                <CardHeader>
                  <CardTitle>Game Console</CardTitle>
                  <CardDescription>Shared panel composition.</CardDescription>
                </CardHeader>
                <CardContent className="text-sm">
                  Queue size: <span className="font-medium">24 players</span>
                </CardContent>
                <CardFooter>
                  <Button size="sm" variant="outline">
                    Inspect
                  </Button>
                </CardFooter>
              </Card>
            </Preview>

            <Preview title="Empty">
              <Empty className="min-h-48 border">
                <EmptyHeader>
                  <EmptyMedia variant="icon">
                    <InboxIcon />
                  </EmptyMedia>
                  <EmptyTitle>No records</EmptyTitle>
                  <EmptyDescription>
                    The state remains visible without adding layout noise.
                  </EmptyDescription>
                </EmptyHeader>
                <EmptyContent>
                  <Button size="sm" variant="outline">
                    Refresh
                  </Button>
                </EmptyContent>
              </Empty>
            </Preview>

            <Preview title="Kbd">
              <KbdGroup>
                <Kbd>⌘</Kbd>
                <Kbd>K</Kbd>
              </KbdGroup>
            </Preview>

            <Preview title="Progress / Skeleton / Spinner">
              <div className="grid gap-4">
                <Progress value={68} />
                <div className="flex items-center gap-3">
                  <Skeleton className="size-10 rounded-md" />
                  <div className="grid flex-1 gap-2">
                    <Skeleton className="h-3 w-2/3" />
                    <Skeleton className="h-3 w-1/2" />
                  </div>
                  <Spinner />
                </div>
              </div>
            </Preview>
          </GallerySection>

          <GallerySection
            id="forms"
            title="Forms"
            description="Inputs, selection controls, form wiring, and field composition."
          >
            <Preview title="Label / Input / Textarea">
              <div className="grid gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="demo-name">Codename</Label>
                  <Input id="demo-name" defaultValue="Pebble Runner" />
                </div>
                <Textarea defaultValue="Meet near the main stage before the round starts." />
              </div>
            </Preview>

            <Preview title="Checkbox / RadioGroup / Switch">
              <div className="grid gap-4">
                <Field orientation="horizontal">
                  <Checkbox id="notifications" defaultChecked />
                  <FieldContent>
                    <FieldLabel htmlFor="notifications">
                      Notifications
                    </FieldLabel>
                    <FieldDescription>
                      Mission updates enabled.
                    </FieldDescription>
                  </FieldContent>
                </Field>
                <RadioGroup defaultValue="team" className="grid gap-2">
                  <Field orientation="horizontal">
                    <RadioGroupItem value="team" id="team" />
                    <FieldLabel htmlFor="team">Team mode</FieldLabel>
                  </Field>
                  <Field orientation="horizontal">
                    <RadioGroupItem value="solo" id="solo" />
                    <FieldLabel htmlFor="solo">Solo mode</FieldLabel>
                  </Field>
                </RadioGroup>
                <Field orientation="horizontal">
                  <Switch id="power" defaultChecked />
                  <FieldLabel htmlFor="power">Power boost</FieldLabel>
                </Field>
              </div>
            </Preview>

            <Preview title="Field">
              <FieldSet>
                <FieldLegend>Player setup</FieldLegend>
                <FieldGroup>
                  <Field>
                    <FieldLabel>Display name</FieldLabel>
                    <Input defaultValue="Yoru" />
                    <FieldDescription>
                      Shown on the leaderboard.
                    </FieldDescription>
                  </Field>
                  <FieldSeparator>Optional</FieldSeparator>
                  <Field data-invalid>
                    <FieldTitle>Team token</FieldTitle>
                    <FieldError>Token is already claimed.</FieldError>
                  </Field>
                </FieldGroup>
              </FieldSet>
            </Preview>

            <Preview title="Form">
              <DemoForm />
            </Preview>

            <Preview title="InputGroup">
              <div className="grid gap-3">
                <InputGroup>
                  <InputGroupAddon>
                    <SearchIcon />
                  </InputGroupAddon>
                  <InputGroupInput placeholder="Search missions" />
                  <InputGroupAddon align="inline-end">
                    <InputGroupButton>
                      <FilterIcon />
                      Filter
                    </InputGroupButton>
                  </InputGroupAddon>
                </InputGroup>
                <InputGroup>
                  <InputGroupAddon align="block-start">
                    <InputGroupText>Notes</InputGroupText>
                  </InputGroupAddon>
                  <InputGroupTextarea defaultValue="Bring badge, water, and a charged phone." />
                </InputGroup>
              </div>
            </Preview>

            <Preview title="InputOTP">
              <InputOTP maxLength={6}>
                <InputOTPGroup>
                  <InputOTPSlot index={0} />
                  <InputOTPSlot index={1} />
                  <InputOTPSlot index={2} />
                </InputOTPGroup>
                <InputOTPSeparator />
                <InputOTPGroup>
                  <InputOTPSlot index={3} />
                  <InputOTPSlot index={4} />
                  <InputOTPSlot index={5} />
                </InputOTPGroup>
              </InputOTP>
            </Preview>

            <Preview title="NativeSelect / Select">
              <div className="grid gap-3 sm:grid-cols-2">
                <NativeSelect defaultValue="alpha">
                  <NativeSelectOptGroup label="Rooms">
                    <NativeSelectOption value="alpha">Alpha</NativeSelectOption>
                    <NativeSelectOption value="beta">Beta</NativeSelectOption>
                  </NativeSelectOptGroup>
                </NativeSelect>
                <Select defaultValue="explore">
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Mode" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>Mode</SelectLabel>
                      <SelectItem value="explore">Explore</SelectItem>
                      <SelectItem value="battle">Battle</SelectItem>
                      <SelectSeparator />
                      <SelectItem value="rest">Rest</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            </Preview>

            <Preview title="Slider">
              <div className="grid gap-3">
                <Slider defaultValue={[64]} max={100} step={1} />
                <div className="text-muted-foreground text-sm">64 points</div>
              </div>
            </Preview>

            <Preview title="Toggle / ToggleGroup">
              <div className="flex flex-wrap gap-3">
                <Toggle variant="outline" defaultPressed aria-label="Star">
                  <StarIcon />
                </Toggle>
                <ToggleGroup
                  type="multiple"
                  variant="outline"
                  defaultValue={["a"]}
                >
                  <ToggleGroupItem value="a" aria-label="A">
                    A
                  </ToggleGroupItem>
                  <ToggleGroupItem value="b" aria-label="B">
                    B
                  </ToggleGroupItem>
                  <ToggleGroupItem value="c" aria-label="C">
                    C
                  </ToggleGroupItem>
                </ToggleGroup>
              </div>
            </Preview>

            <Preview title="ButtonGroup">
              <ButtonGroup>
                <Button variant="outline">
                  <PlusIcon />
                  Add
                </Button>
                <ButtonGroupSeparator />
                <ButtonGroupText>12</ButtonGroupText>
                <Button variant="outline" size="icon" aria-label="More">
                  <ChevronDownIcon />
                </Button>
              </ButtonGroup>
            </Preview>
          </GallerySection>

          <GallerySection
            id="overlays"
            title="Overlays"
            description="Dialogs, menus, popovers, drawers, and command surfaces."
          >
            <Preview title="Dialog">
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline">Open dialog</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm round</DialogTitle>
                    <DialogDescription>
                      This preview uses the shared dialog primitives.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter showCloseButton>
                    <Button>Start</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </Preview>

            <Preview title="AlertDialog">
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">
                    <Trash2Icon />
                    Reset
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogMedia>
                      <ShieldAlertIcon />
                    </AlertDialogMedia>
                    <AlertDialogTitle>Reset scoreboard?</AlertDialogTitle>
                    <AlertDialogDescription>
                      Scores will return to the demo baseline.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction variant="destructive">
                      Reset
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </Preview>

            <Preview title="Sheet / Drawer">
              <div className="flex flex-wrap gap-2">
                <Sheet>
                  <SheetTrigger asChild>
                    <Button variant="outline">Open sheet</Button>
                  </SheetTrigger>
                  <SheetContent>
                    <SheetHeader>
                      <SheetTitle>Mission queue</SheetTitle>
                      <SheetDescription>
                        Three tasks are ready.
                      </SheetDescription>
                    </SheetHeader>
                    <SheetFooter>
                      <SheetClose asChild>
                        <Button>Close</Button>
                      </SheetClose>
                    </SheetFooter>
                  </SheetContent>
                </Sheet>
                <Drawer>
                  <DrawerTrigger asChild>
                    <Button variant="outline">Open drawer</Button>
                  </DrawerTrigger>
                  <DrawerContent>
                    <DrawerHeader>
                      <DrawerTitle>Round details</DrawerTitle>
                      <DrawerDescription>
                        Bottom drawer preview.
                      </DrawerDescription>
                    </DrawerHeader>
                    <DrawerFooter>
                      <DrawerClose asChild>
                        <Button>Done</Button>
                      </DrawerClose>
                    </DrawerFooter>
                  </DrawerContent>
                </Drawer>
              </div>
            </Preview>

            <Preview title="DropdownMenu">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline">
                    Actions
                    <ChevronDownIcon />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start">
                  <DropdownMenuLabel>Round</DropdownMenuLabel>
                  <DropdownMenuItem>
                    <CopyIcon />
                    Duplicate
                    <DropdownMenuShortcut>⌘D</DropdownMenuShortcut>
                  </DropdownMenuItem>
                  <DropdownMenuCheckboxItem checked>
                    Auto publish
                  </DropdownMenuCheckboxItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem variant="destructive">
                    <Trash2Icon />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </Preview>

            <Preview title="ContextMenu">
              <ContextMenu>
                <ContextMenuTrigger className="bg-muted/50 flex h-28 items-center justify-center rounded-md border border-dashed text-sm">
                  Right click area
                </ContextMenuTrigger>
                <ContextMenuContent>
                  <ContextMenuLabel>Board</ContextMenuLabel>
                  <ContextMenuItem>
                    <CopyIcon />
                    Copy
                    <ContextMenuShortcut>⌘C</ContextMenuShortcut>
                  </ContextMenuItem>
                  <ContextMenuCheckboxItem checked>
                    Show labels
                  </ContextMenuCheckboxItem>
                  <ContextMenuSeparator />
                  <ContextMenuItem variant="destructive">
                    Remove
                  </ContextMenuItem>
                </ContextMenuContent>
              </ContextMenu>
            </Preview>

            <Preview title="HoverCard / Popover / Tooltip">
              <div className="flex flex-wrap gap-2">
                <HoverCard>
                  <HoverCardTrigger asChild>
                    <Button variant="outline">Hover card</Button>
                  </HoverCardTrigger>
                  <HoverCardContent>
                    <div className="space-y-1">
                      <div className="font-medium">Team Alpha</div>
                      <p className="text-muted-foreground text-sm">
                        Current rank: 2
                      </p>
                    </div>
                  </HoverCardContent>
                </HoverCard>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button variant="outline">Popover</Button>
                  </PopoverTrigger>
                  <PopoverContent>
                    <PopoverHeader>
                      <PopoverTitle>Quick filter</PopoverTitle>
                      <PopoverDescription>
                        Apply a temporary board filter.
                      </PopoverDescription>
                    </PopoverHeader>
                  </PopoverContent>
                </Popover>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="outline" size="icon" aria-label="Help">
                      <CircleHelpIcon />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Tooltip content</TooltipContent>
                </Tooltip>
              </div>
            </Preview>

            <Preview title="Command">
              <Command className="rounded-md border">
                <CommandInput placeholder="Search command" />
                <CommandList>
                  <CommandEmpty>No command found.</CommandEmpty>
                  <CommandGroup heading="Actions">
                    <CommandItem>
                      <SparklesIcon />
                      Start challenge
                      <CommandShortcut>⌘S</CommandShortcut>
                    </CommandItem>
                    <CommandItem>
                      <MailIcon />
                      Message team
                    </CommandItem>
                  </CommandGroup>
                  <CommandSeparator />
                  <CommandGroup heading="System">
                    <CommandItem>
                      <SettingsIcon />
                      Settings
                    </CommandItem>
                  </CommandGroup>
                </CommandList>
              </Command>
            </Preview>

            <Preview title="Combobox">
              <Combobox items={["Mission board", "Item shop", "Leaderboard"]}>
                <ComboboxInput placeholder="Search destination" />
                <ComboboxContent>
                  <ComboboxList>
                    <ComboboxItem value="Mission board">
                      Mission board
                    </ComboboxItem>
                    <ComboboxItem value="Item shop">Item shop</ComboboxItem>
                    <ComboboxItem value="Leaderboard">Leaderboard</ComboboxItem>
                    <ComboboxEmpty>No destination.</ComboboxEmpty>
                  </ComboboxList>
                </ComboboxContent>
              </Combobox>
            </Preview>

            <Preview title="Sonner">
              <Button
                variant="outline"
                onClick={() => toast.success("Toast queued")}
              >
                <BellIcon />
                Toast
              </Button>
            </Preview>
          </GallerySection>

          <GallerySection
            id="data"
            title="Data"
            description="Tables, charts, calendar, lists, tabs, and scrollable content."
          >
            <Preview title="Calendar">
              <Calendar mode="single" className="rounded-md border" />
            </Preview>

            <Preview title="Chart">
              <ChartContainer config={chartConfig} className="min-h-56">
                <RechartsAreaChart data={chartData}>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="day"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                  />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Area
                    dataKey="score"
                    type="natural"
                    fill="var(--color-score)"
                    fillOpacity={0.22}
                    stroke="var(--color-score)"
                  />
                </RechartsAreaChart>
              </ChartContainer>
            </Preview>

            <Preview title="Table">
              <Table>
                <TableCaption>Demo leaderboard</TableCaption>
                <TableHeader>
                  <TableRow>
                    <TableHead>Team</TableHead>
                    <TableHead>Round</TableHead>
                    <TableHead className="text-right">Score</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow>
                    <TableCell>Alpha</TableCell>
                    <TableCell>Quiz</TableCell>
                    <TableCell className="text-right">128</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Beta</TableCell>
                    <TableCell>Fusion</TableCell>
                    <TableCell className="text-right">96</TableCell>
                  </TableRow>
                </TableBody>
                <TableFooter>
                  <TableRow>
                    <TableCell colSpan={2}>Total</TableCell>
                    <TableCell className="text-right">224</TableCell>
                  </TableRow>
                </TableFooter>
              </Table>
            </Preview>

            <Preview title="Tabs">
              <Tabs defaultValue="overview" className="w-full">
                <TabsList>
                  <TabsTrigger value="overview">Overview</TabsTrigger>
                  <TabsTrigger value="activity">Activity</TabsTrigger>
                </TabsList>
                <TabsContent value="overview" className="text-sm">
                  Active missions and active players.
                </TabsContent>
                <TabsContent value="activity" className="text-sm">
                  Recent joins, answers, and fusions.
                </TabsContent>
              </Tabs>
            </Preview>

            <Preview title="Item">
              <ItemGroup className="rounded-md border">
                <Item>
                  <ItemMedia variant="icon">
                    <ClipboardListIcon />
                  </ItemMedia>
                  <ItemContent>
                    <ItemTitle>Mission checkpoint</ItemTitle>
                    <ItemDescription>
                      Three teams have checked in.
                    </ItemDescription>
                  </ItemContent>
                  <ItemActions>
                    <Button size="sm" variant="outline">
                      View
                    </Button>
                  </ItemActions>
                  <ItemFooter>
                    <Badge variant="secondary">Live</Badge>
                  </ItemFooter>
                </Item>
                <ItemSeparator />
                <Item>
                  <ItemHeader>
                    <ItemTitle>Fusion recipe</ItemTitle>
                    <Badge variant="outline">New</Badge>
                  </ItemHeader>
                  <ItemDescription>
                    Reusable item layout primitives.
                  </ItemDescription>
                </Item>
              </ItemGroup>
            </Preview>

            <Preview title="ScrollArea">
              <ScrollArea className="h-44 rounded-md border">
                <div className="grid gap-2 p-4">
                  {Array.from({ length: 12 }, (_, index) => (
                    <div
                      key={index}
                      className="bg-muted/50 rounded-md px-3 py-2 text-sm"
                    >
                      Event log #{index + 1}
                    </div>
                  ))}
                </div>
                <ScrollBar />
              </ScrollArea>
            </Preview>
          </GallerySection>

          <GallerySection
            id="navigation"
            title="Navigation"
            description="Breadcrumbs, menus, pagination, and navigation menu primitives."
          >
            <Preview title="Breadcrumb">
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem>
                    <BreadcrumbLink href="/">Camp</BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator />
                  <BreadcrumbItem>
                    <BreadcrumbLink href="/component-gallery">
                      UI
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator />
                  <BreadcrumbItem>
                    <BreadcrumbPage>Gallery</BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </Preview>

            <Preview title="Menubar">
              <Menubar>
                <MenubarMenu>
                  <MenubarTrigger>Game</MenubarTrigger>
                  <MenubarContent>
                    <MenubarItem>
                      New round
                      <MenubarShortcut>⌘N</MenubarShortcut>
                    </MenubarItem>
                    <MenubarItem>Export</MenubarItem>
                    <MenubarSeparator />
                    <MenubarItem>Settings</MenubarItem>
                  </MenubarContent>
                </MenubarMenu>
                <MenubarMenu>
                  <MenubarTrigger>View</MenubarTrigger>
                  <MenubarContent>
                    <MenubarItem>Leaderboard</MenubarItem>
                    <MenubarItem>Inventory</MenubarItem>
                  </MenubarContent>
                </MenubarMenu>
              </Menubar>
            </Preview>

            <Preview title="NavigationMenu">
              <NavigationMenu viewport={false}>
                <NavigationMenuList>
                  <NavigationMenuItem>
                    <NavigationMenuTrigger>Modules</NavigationMenuTrigger>
                    <NavigationMenuContent>
                      <div className="grid w-64 gap-1">
                        <NavigationMenuLink href="/" active>
                          <HomeIcon />
                          Home
                        </NavigationMenuLink>
                        <NavigationMenuLink href="/component-gallery">
                          <LayersIcon />
                          Components
                        </NavigationMenuLink>
                      </div>
                    </NavigationMenuContent>
                  </NavigationMenuItem>
                  <NavigationMenuItem>
                    <NavigationMenuLink href="/" className="px-4 py-2">
                      Status
                    </NavigationMenuLink>
                  </NavigationMenuItem>
                </NavigationMenuList>
              </NavigationMenu>
            </Preview>

            <Preview title="Pagination">
              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious href="#" />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#" isActive>
                      1
                    </PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#">2</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationEllipsis />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationNext href="#" />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </Preview>
          </GallerySection>

          <GallerySection
            id="layout"
            title="Layout"
            description="Collapsible, carousel, direction, resizable panels, separator, and sidebar."
          >
            <Preview title="Collapsible">
              <Collapsible className="grid gap-2">
                <CollapsibleTrigger asChild>
                  <Button variant="outline">
                    <ChevronsUpDownIcon />
                    Toggle details
                  </Button>
                </CollapsibleTrigger>
                <CollapsibleContent className="rounded-md border p-3 text-sm">
                  Extra round metadata can live behind this control.
                </CollapsibleContent>
              </Collapsible>
            </Preview>

            <Preview title="Carousel">
              <Carousel className="mx-auto w-full max-w-sm">
                <CarouselContent>
                  {[1, 2, 3].map((item) => (
                    <CarouselItem key={item}>
                      <div className="bg-muted flex aspect-video items-center justify-center rounded-md border text-2xl font-semibold">
                        {item}
                      </div>
                    </CarouselItem>
                  ))}
                </CarouselContent>
                <CarouselPrevious className="left-2" />
                <CarouselNext className="right-2" />
              </Carousel>
            </Preview>

            <Preview title="Direction">
              <DirectionProvider dir="rtl" direction="rtl">
                <div className="rounded-md border p-3 text-sm">
                  RTL preview · <span className="font-medium">مرحبا</span>
                </div>
              </DirectionProvider>
            </Preview>

            <Preview title="Resizable">
              <ResizablePanelGroup
                orientation="horizontal"
                className="min-h-40 rounded-md border"
              >
                <ResizablePanel defaultSize={45}>
                  <div className="flex h-full items-center justify-center text-sm">
                    Left
                  </div>
                </ResizablePanel>
                <ResizableHandle withHandle />
                <ResizablePanel defaultSize={55}>
                  <div className="flex h-full items-center justify-center text-sm">
                    Right
                  </div>
                </ResizablePanel>
              </ResizablePanelGroup>
            </Preview>

            <Preview title="Separator">
              <div className="space-y-3">
                <div className="flex h-5 items-center gap-3 text-sm">
                  <span>Alpha</span>
                  <Separator orientation="vertical" />
                  <span>Beta</span>
                </div>
                <Separator />
              </div>
            </Preview>

            <Preview title="Sidebar">
              <div className="overflow-hidden rounded-md border">
                <SidebarProvider defaultOpen>
                  <Sidebar collapsible="none" className="min-h-64">
                    <SidebarHeader>
                      <div className="flex items-center gap-2 px-2 py-1 text-sm font-medium">
                        <LayoutDashboardIcon className="size-4" />
                        Console
                      </div>
                    </SidebarHeader>
                    <SidebarContent>
                      <SidebarGroup>
                        <SidebarGroupLabel>Menu</SidebarGroupLabel>
                        <SidebarGroupContent>
                          <SidebarMenu>
                            <SidebarMenuItem>
                              <SidebarMenuButton isActive>
                                <HomeIcon />
                                <span>Home</span>
                              </SidebarMenuButton>
                            </SidebarMenuItem>
                            <SidebarMenuItem>
                              <SidebarMenuButton>
                                <LineChartIcon />
                                <span>Reports</span>
                              </SidebarMenuButton>
                            </SidebarMenuItem>
                            <SidebarMenuItem>
                              <SidebarMenuSkeleton showIcon />
                            </SidebarMenuItem>
                          </SidebarMenu>
                        </SidebarGroupContent>
                      </SidebarGroup>
                    </SidebarContent>
                  </Sidebar>
                  <SidebarInset className="min-h-64 p-4">
                    <div className="flex items-center gap-2">
                      <SidebarTrigger />
                      <span className="text-sm font-medium">Inset</span>
                    </div>
                  </SidebarInset>
                </SidebarProvider>
              </div>
            </Preview>
          </GallerySection>

          <GallerySection
            id="project-ui"
            title="Project UI"
            description="Local shared wrappers built on top of the primitives."
          >
            <Preview title="IconBadge / StatusBadge / StatusDot">
              <div className="flex flex-wrap items-center gap-2">
                <IconBadge
                  label="Explore"
                  tone="explore"
                  icon={<SparklesIcon className="size-4" />}
                />
                <StatusBadge tone="success">Online</StatusBadge>
                <span className="inline-flex items-center gap-2 text-sm">
                  <StatusDot tone="magic" />
                  Magic
                </span>
              </div>
            </Preview>

            <Preview title="InfoTile / MetricCard">
              <div className="grid gap-3 sm:grid-cols-2">
                <InfoTile
                  label="Open matches"
                  value="8"
                  icon={RadioTowerIcon}
                />
                <MetricCard
                  label="Knowledge battles"
                  value="5"
                  tone="info"
                  icon={BookOpenIcon}
                />
              </div>
            </Preview>
          </GallerySection>
        </div>
      </main>
    </TooltipProvider>
  )
}

function GallerySection({
  id,
  title,
  description,
  children,
}: {
  id: string
  title: string
  description: string
  children: ReactNode
}) {
  return (
    <section id={id} className="scroll-mt-6 space-y-4 border-t-2 pt-6 pb-2">
      <div className="space-y-1">
        <p className="text-muted-foreground text-xs font-black tracking-normal uppercase">
          Component Set
        </p>
        <h2 className="text-2xl font-black">{title}</h2>
        <p className="text-muted-foreground text-sm">{description}</p>
      </div>
      <div className="grid gap-4 lg:grid-cols-2">{children}</div>
    </section>
  )
}

function Preview({ title, children }: { title: string; children: ReactNode }) {
  return (
    <article className="bg-card text-card-foreground border-ink rounded-[1.375rem] border-2 p-4 shadow-[4px_4px_0_rgba(23,35,58,0.12)]">
      <div className="mb-4 flex items-center justify-between gap-3">
        <h3 className="text-sm font-black">{title}</h3>
        <Badge variant="outline" className="border-border border-2 font-black">
          {title}
        </Badge>
      </div>
      <div className="min-h-20">{children}</div>
    </article>
  )
}

function DemoForm() {
  const form = useForm<DemoFormValues>({
    defaultValues: {
      codename: "Pebble Runner",
    },
  })

  return (
    <Form {...form}>
      <form
        className="grid gap-3"
        onSubmit={(event) => {
          event.preventDefault()
        }}
      >
        <FormField
          control={form.control}
          name="codename"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Codename</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormDescription>Controlled by react-hook-form.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="w-fit">
          Submit
        </Button>
      </form>
    </Form>
  )
}
