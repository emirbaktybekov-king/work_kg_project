"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api, Job, Stats, User, Resume } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Briefcase,
  Users,
  Plus,
  Pencil,
  Trash2,
  LogOut,
  TrendingUp,
  Menu,
  X,
  ChevronLeft,
  ChevronRight,
  LayoutDashboard,
  Phone,
  MapPin,
  Calendar,
  MessageCircle,
  ExternalLink,
  User as UserIcon,
  FileText,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const categories: Record<string, string[]> = {
  "Строительство": ["Каменщик", "Кладка", "Электрик", "Сантехник", "Сварщик", "Отделочник", "Плиточник", "Фасадчик", "Монолитчик", "Разнорабочий"],
  "Общепит": ["Повар", "Официант", "Бармен", "Посудомойщик", "Администратор", "Кассир"],
  "Швейный цех": ["Швея", "Закройщик", "Упаковщик", "Технолог", "Контролер качества"],
  "IT": ["Программист", "Дизайнер", "Тестировщик", "Системный администратор"],
  "Продажи": ["Продавец", "Менеджер", "Консультант", "Кассир"],
  "Транспорт": ["Водитель", "Курьер", "Экспедитор", "Диспетчер"],
};

const cities = ["Бишкек", "Ош", "Талас", "Нарын", "Каракол", "Жалал-Абад", "Чолпон-Ата"];

export default function CRMPage() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<"dashboard" | "resumes" | "jobs" | "users">("dashboard");
  const [jobs, setJobs] = useState<Job[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [resumes, setResumes] = useState<Resume[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState<string>("");

  // Sidebar state
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [isDesktop, setIsDesktop] = useState(true);

  // User details modal
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [userModalOpen, setUserModalOpen] = useState(false);

  // Handle window resize
  useEffect(() => {
    const handleResize = () => {
      setIsDesktop(window.innerWidth >= 1024);
    };
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Job dialog state
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingJob, setEditingJob] = useState<Job | null>(null);
  const [jobForm, setJobForm] = useState({
    title: "",
    description: "",
    category: "",
    subcategory: "",
    city: "",
    salary: "",
    phone: "",
    company: "",
    is_active: true,
  });

  useEffect(() => {
    const token = api.getToken();
    if (!token) {
      router.push("/auth/login");
      return;
    }

    const role = localStorage.getItem("role");
    setUserRole(role || "");

    loadData();
  }, [router]);

  const loadData = async () => {
    try {
      const [jobsData, usersData, resumesData, statsData] = await Promise.all([
        api.getJobs(),
        api.getUsers(),
        api.getResumes(),
        api.getStats(),
      ]);
      setJobs(jobsData || []);
      setUsers(usersData || []);
      setResumes(resumesData || []);
      setStats(statsData);
    } catch (error) {
      console.error("Failed to load data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    api.logout();
    router.push("/auth/login");
  };

  const openCreateDialog = () => {
    setEditingJob(null);
    setJobForm({
      title: "",
      description: "",
      category: "",
      subcategory: "",
      city: "",
      salary: "",
      phone: "",
      company: "",
      is_active: true,
    });
    setIsDialogOpen(true);
  };

  const openEditDialog = (job: Job) => {
    setEditingJob(job);
    setJobForm({
      title: job.title,
      description: job.description,
      category: job.category,
      subcategory: job.subcategory,
      city: job.city,
      salary: job.salary,
      phone: job.phone,
      company: job.company,
      is_active: job.is_active,
    });
    setIsDialogOpen(true);
  };

  const handleSaveJob = async () => {
    try {
      if (editingJob) {
        await api.updateJob(editingJob.id, jobForm);
      } else {
        await api.createJob(jobForm);
      }
      setIsDialogOpen(false);
      loadData();
    } catch (error) {
      console.error("Failed to save job:", error);
    }
  };

  const handleDeleteJob = async (id: number) => {
    if (!confirm("Вы уверены, что хотите удалить эту вакансию?")) return;
    try {
      await api.deleteJob(id);
      loadData();
    } catch (error) {
      console.error("Failed to delete job:", error);
    }
  };

  const openUserModal = (user: User) => {
    setSelectedUser(user);
    setUserModalOpen(true);
  };

  const menuItems = [
    { id: "dashboard", label: "Панель управления", icon: LayoutDashboard },
    { id: "resumes", label: "Анкеты", icon: FileText },
    { id: "jobs", label: "Вакансии", icon: Briefcase },
    { id: "users", label: "Пользователи", icon: Users },
  ];

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="text-lg font-medium text-gray-600"
        >
          Загрузка...
        </motion.div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 flex">
      {/* Mobile menu overlay */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-40 lg:hidden"
            onClick={() => setMobileMenuOpen(false)}
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <AnimatePresence mode="wait">
        <motion.aside
          initial={false}
          animate={{
            width: sidebarOpen ? 256 : 80,
            x: mobileMenuOpen ? 0 : (!isDesktop ? -256 : 0),
          }}
          transition={{ duration: 0.3, ease: "easeInOut" }}
          className={`fixed left-0 top-0 h-full bg-white border-r border-gray-200 z-50 flex flex-col
            ${mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}`}
        >
          {/* Logo */}
          <div className="h-16 flex items-center justify-between px-4 border-b border-gray-200">
            <motion.div
              initial={false}
              animate={{ opacity: sidebarOpen ? 1 : 0, width: sidebarOpen ? "auto" : 0 }}
              className="flex items-center gap-3 overflow-hidden"
            >
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-600 to-blue-700 flex-shrink-0">
                <span className="text-sm font-bold text-white">WK</span>
              </div>
              <div className="whitespace-nowrap">
                <h1 className="text-lg font-bold text-gray-900">WorkKG</h1>
                <p className="text-xs text-gray-500">Админ панель</p>
              </div>
            </motion.div>

            {/* Close button for mobile */}
            <button
              onClick={() => setMobileMenuOpen(false)}
              className="lg:hidden p-2 rounded-lg hover:bg-gray-100"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 p-4 space-y-2">
            {menuItems.map((item) => {
              const Icon = item.icon;
              const isActive = activeTab === item.id;
              return (
                <motion.button
                  key={item.id}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={() => {
                    setActiveTab(item.id as typeof activeTab);
                    setMobileMenuOpen(false);
                  }}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all
                    ${isActive
                      ? 'bg-blue-50 text-blue-600 font-medium'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                    }`}
                >
                  <Icon className={`h-5 w-5 flex-shrink-0 ${isActive ? 'text-blue-600' : 'text-gray-400'}`} />
                  <motion.span
                    initial={false}
                    animate={{ opacity: sidebarOpen ? 1 : 0, width: sidebarOpen ? "auto" : 0 }}
                    className="whitespace-nowrap overflow-hidden"
                  >
                    {item.label}
                  </motion.span>
                </motion.button>
              );
            })}
          </nav>

          {/* Collapse toggle */}
          <div className="p-4 border-t border-gray-200 hidden lg:block">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 rounded-xl text-gray-500 hover:bg-gray-100 transition-colors"
            >
              {sidebarOpen ? (
                <>
                  <ChevronLeft className="h-5 w-5" />
                  <span className="text-sm">Свернуть</span>
                </>
              ) : (
                <ChevronRight className="h-5 w-5" />
              )}
            </button>
          </div>
        </motion.aside>
      </AnimatePresence>

      {/* Main content wrapper */}
      <motion.div
        initial={false}
        animate={{ marginLeft: isDesktop ? (sidebarOpen ? 256 : 80) : 0 }}
        transition={{ duration: 0.3, ease: "easeInOut" }}
        className="flex-1 flex flex-col min-h-screen"
      >
        {/* Header */}
        <header className="h-16 bg-white border-b border-gray-200 sticky top-0 z-30 flex items-center justify-between px-4 lg:px-6">
          <div className="flex items-center gap-4">
            {/* Mobile menu button */}
            <button
              onClick={() => setMobileMenuOpen(true)}
              className="lg:hidden p-2 rounded-lg hover:bg-gray-100"
            >
              <Menu className="h-6 w-6" />
            </button>
            <h2 className="text-lg font-semibold text-gray-900">
              {activeTab === "dashboard" && "Панель управления"}
              {activeTab === "resumes" && "Анкеты"}
              {activeTab === "jobs" && "Вакансии"}
              {activeTab === "users" && "Пользователи"}
            </h2>
          </div>

          <div className="flex items-center gap-3">
            <Badge variant="secondary" className="hidden sm:inline-flex">
              {userRole || "admin"}
            </Badge>
            <Button variant="outline" size="sm" onClick={handleLogout}>
              <LogOut className="h-4 w-4 mr-2" />
              <span className="hidden sm:inline">Выйти</span>
            </Button>
          </div>
        </header>

        {/* Main content */}
        <main className="flex-1 p-4 lg:p-6 overflow-auto">
          <AnimatePresence mode="wait">
            {activeTab === "dashboard" && (
              <motion.div
                key="dashboard"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-6"
              >
                {/* Stats cards */}
                <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
                  {[
                    { title: "Анкеты", value: stats?.total_resumes || 0, sub: `${stats?.today_resumes || 0} сегодня`, icon: FileText, color: "indigo" },
                    { title: "Вакансии", value: stats?.total_jobs || 0, sub: `${stats?.active_jobs || 0} активных`, icon: Briefcase, color: "blue" },
                    { title: "Пользователи", value: stats?.total_users || 0, sub: `${stats?.today_users || 0} сегодня`, icon: Users, color: "purple" },
                    { title: "Статус бота", value: "Онлайн", sub: "Telegram бот", icon: MessageCircle, color: "emerald" },
                  ].map((stat, idx) => (
                    <motion.div
                      key={stat.title}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: idx * 0.1 }}
                    >
                      <Card className="hover:shadow-md transition-shadow">
                        <CardHeader className="flex flex-row items-center justify-between pb-2">
                          <CardTitle className="text-sm font-medium text-gray-500">{stat.title}</CardTitle>
                          <stat.icon className={`h-4 w-4 text-${stat.color}-500`} />
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold">{stat.value}</div>
                          <p className="text-xs text-gray-400">{stat.sub}</p>
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
                </div>

                {/* Recent jobs */}
                <Card>
                  <CardHeader>
                    <CardTitle>Последние вакансии</CardTitle>
                  </CardHeader>
                  <CardContent className="p-0">
                    {/* Desktop table */}
                    <div className="hidden md:block">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>Название</TableHead>
                            <TableHead>Категория</TableHead>
                            <TableHead>Город</TableHead>
                            <TableHead>Источник</TableHead>
                            <TableHead>Статус</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {jobs.slice(0, 5).map((job) => (
                            <TableRow key={job.id}>
                              <TableCell className="font-medium">{job.title}</TableCell>
                              <TableCell>{job.category}</TableCell>
                              <TableCell>{job.city}</TableCell>
                              <TableCell>
                                <Badge variant="outline">{job.source}</Badge>
                              </TableCell>
                              <TableCell>
                                <Badge variant={job.is_active ? "default" : "secondary"}>
                                  {job.is_active ? "Активна" : "Неактивна"}
                                </Badge>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>

                    {/* Mobile cards */}
                    <div className="md:hidden p-4 space-y-3">
                      {jobs.slice(0, 5).map((job) => (
                        <motion.div
                          key={job.id}
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                          className="p-4 bg-gray-50 rounded-xl space-y-2"
                        >
                          <div className="font-medium">{job.title}</div>
                          <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                            <span>{job.category}</span>
                            <span>•</span>
                            <span>{job.city}</span>
                          </div>
                          <div className="flex gap-2">
                            <Badge variant="outline">{job.source}</Badge>
                            <Badge variant={job.is_active ? "default" : "secondary"}>
                              {job.is_active ? "Активна" : "Неактивна"}
                            </Badge>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            )}

            {activeTab === "resumes" && (
              <motion.div
                key="resumes"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-6"
              >
                <h2 className="text-2xl font-bold">Анкеты соискателей</h2>

                {/* Desktop table */}
                <Card className="hidden lg:block">
                  <CardContent className="p-0">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>Telegram</TableHead>
                          <TableHead>Имя</TableHead>
                          <TableHead>Телефон</TableHead>
                          <TableHead>Город</TableHead>
                          <TableHead>Специальность</TableHead>
                          <TableHead>Опыт</TableHead>
                          <TableHead>Обновлено</TableHead>
                          <TableHead>Действия</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {resumes.map((resume) => (
                          <TableRow key={resume.id}>
                            <TableCell>{resume.id}</TableCell>
                            <TableCell>
                              {resume.username ? (
                                <a
                                  href={`https://t.me/${resume.username}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="text-blue-600 hover:underline"
                                >
                                  @{resume.username}
                                </a>
                              ) : (
                                <span className="text-gray-400">ID: {resume.telegram_id}</span>
                              )}
                            </TableCell>
                            <TableCell className="font-medium">{resume.name || "-"}</TableCell>
                            <TableCell>{resume.phone || "-"}</TableCell>
                            <TableCell>{resume.city || "-"}</TableCell>
                            <TableCell>{resume.specialty || "-"}</TableCell>
                            <TableCell className="max-w-[200px] truncate" title={resume.experience}>
                              {resume.experience || "-"}
                            </TableCell>
                            <TableCell>
                              {new Date(resume.updated_at).toLocaleDateString("ru-RU")}
                            </TableCell>
                            <TableCell>
                              {resume.username && (
                                <a
                                  href={`https://t.me/${resume.username}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                >
                                  <Button size="sm" variant="outline">
                                    <ExternalLink className="h-3 w-3 mr-1" /> Telegram
                                  </Button>
                                </a>
                              )}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>

                {/* Mobile/Tablet cards grid */}
                <div className="lg:hidden grid gap-4 sm:grid-cols-2">
                  {resumes.map((resume, idx) => (
                    <motion.div
                      key={resume.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: idx * 0.05 }}
                    >
                      <Card className="hover:shadow-md transition-shadow">
                        <CardContent className="p-4 space-y-3">
                          <div className="flex items-center gap-3">
                            <div className="h-12 w-12 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold">
                              {(resume.name?.[0] || resume.username?.[0] || "A").toUpperCase()}
                            </div>
                            <div>
                              <h3 className="font-medium">{resume.name || "Без имени"}</h3>
                              {resume.username ? (
                                <a
                                  href={`https://t.me/${resume.username}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="text-sm text-blue-600 hover:underline"
                                >
                                  @{resume.username}
                                </a>
                              ) : (
                                <p className="text-sm text-gray-500">ID: {resume.telegram_id}</p>
                              )}
                            </div>
                          </div>

                          <div className="space-y-1 text-sm">
                            {resume.phone && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <Phone className="h-3 w-3" /> {resume.phone}
                              </div>
                            )}
                            {resume.city && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <MapPin className="h-3 w-3" /> {resume.city}
                              </div>
                            )}
                            {resume.specialty && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <Briefcase className="h-3 w-3" /> {resume.specialty}
                              </div>
                            )}
                          </div>

                          {resume.experience && (
                            <div className="text-sm text-gray-600 bg-gray-50 p-2 rounded">
                              <span className="font-medium">Опыт:</span> {resume.experience}
                            </div>
                          )}

                          <div className="flex items-center justify-between pt-2">
                            <div className="flex items-center gap-2 text-xs text-gray-400">
                              <Calendar className="h-3 w-3" />
                              {new Date(resume.updated_at).toLocaleDateString("ru-RU")}
                            </div>
                            {resume.username && (
                              <a
                                href={`https://t.me/${resume.username}`}
                                target="_blank"
                                rel="noopener noreferrer"
                              >
                                <Button size="sm" variant="outline">
                                  <ExternalLink className="h-3 w-3 mr-1" /> Telegram
                                </Button>
                              </a>
                            )}
                          </div>
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
                </div>

                {resumes.length === 0 && (
                  <Card>
                    <CardContent className="p-8 text-center text-gray-500">
                      <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>Анкеты пока не заполнены</p>
                      <p className="text-sm">Когда пользователи заполнят анкеты в боте, они появятся здесь</p>
                    </CardContent>
                  </Card>
                )}
              </motion.div>
            )}

            {activeTab === "jobs" && (
              <motion.div
                key="jobs"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-6"
              >
                <div className="flex items-center justify-between">
                  <h2 className="text-2xl font-bold">Вакансии</h2>
                  <Button onClick={openCreateDialog}>
                    <Plus className="h-4 w-4 mr-2" />
                    Добавить
                  </Button>
                </div>

                {/* Desktop table */}
                <Card className="hidden md:block">
                  <CardContent className="p-0">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>Название</TableHead>
                          <TableHead>Категория</TableHead>
                          <TableHead>Город</TableHead>
                          <TableHead>Зарплата</TableHead>
                          <TableHead>Телефон</TableHead>
                          <TableHead>Источник</TableHead>
                          <TableHead>Статус</TableHead>
                          <TableHead className="text-right">Действия</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {jobs.map((job) => (
                          <TableRow key={job.id}>
                            <TableCell>{job.id}</TableCell>
                            <TableCell className="font-medium max-w-[200px] truncate">
                              {job.title}
                            </TableCell>
                            <TableCell>
                              <div className="text-sm">
                                {job.category}
                                <br />
                                <span className="text-gray-400">{job.subcategory}</span>
                              </div>
                            </TableCell>
                            <TableCell>{job.city}</TableCell>
                            <TableCell>{job.salary}</TableCell>
                            <TableCell>{job.phone}</TableCell>
                            <TableCell>
                              <Badge variant="outline">{job.source}</Badge>
                            </TableCell>
                            <TableCell>
                              <Badge variant={job.is_active ? "default" : "secondary"}>
                                {job.is_active ? "Активна" : "Неактивна"}
                              </Badge>
                            </TableCell>
                            <TableCell className="text-right">
                              <div className="flex justify-end gap-2">
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => openEditDialog(job)}
                                >
                                  <Pencil className="h-4 w-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleDeleteJob(job.id)}
                                >
                                  <Trash2 className="h-4 w-4 text-red-500" />
                                </Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>

                {/* Mobile cards grid */}
                <div className="md:hidden grid gap-4">
                  {jobs.map((job, idx) => (
                    <motion.div
                      key={job.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: idx * 0.05 }}
                    >
                      <Card>
                        <CardContent className="p-4 space-y-3">
                          <div className="flex items-start justify-between">
                            <div>
                              <h3 className="font-medium">{job.title}</h3>
                              <p className="text-sm text-gray-500">{job.category} / {job.subcategory}</p>
                            </div>
                            <Badge variant={job.is_active ? "default" : "secondary"}>
                              {job.is_active ? "Активна" : "Неактивна"}
                            </Badge>
                          </div>

                          <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                            <span className="flex items-center gap-1">
                              <MapPin className="h-3 w-3" /> {job.city}
                            </span>
                            {job.salary && (
                              <span>💰 {job.salary}</span>
                            )}
                          </div>

                          {job.phone && (
                            <div className="flex items-center gap-1 text-sm text-gray-500">
                              <Phone className="h-3 w-3" /> {job.phone}
                            </div>
                          )}

                          <div className="flex gap-2 pt-2">
                            <Button size="sm" variant="outline" onClick={() => openEditDialog(job)}>
                              <Pencil className="h-3 w-3 mr-1" /> Изменить
                            </Button>
                            <Button size="sm" variant="outline" onClick={() => handleDeleteJob(job.id)}>
                              <Trash2 className="h-3 w-3 mr-1 text-red-500" /> Удалить
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
                </div>
              </motion.div>
            )}

            {activeTab === "users" && (
              <motion.div
                key="users"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-6"
              >
                <h2 className="text-2xl font-bold">Пользователи</h2>

                {/* Desktop table */}
                <Card className="hidden lg:block">
                  <CardContent className="p-0">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>Telegram ID</TableHead>
                          <TableHead>Username</TableHead>
                          <TableHead>Имя</TableHead>
                          <TableHead>Телефон</TableHead>
                          <TableHead>Город</TableHead>
                          <TableHead>Специальность</TableHead>
                          <TableHead>Опыт</TableHead>
                          <TableHead>Регистрация</TableHead>
                          <TableHead>Действия</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {users.map((user) => (
                          <TableRow key={user.id} className="cursor-pointer hover:bg-gray-50" onClick={() => openUserModal(user)}>
                            <TableCell>{user.id}</TableCell>
                            <TableCell>{user.telegram_id}</TableCell>
                            <TableCell>@{user.username || "-"}</TableCell>
                            <TableCell>
                              {user.first_name} {user.last_name}
                            </TableCell>
                            <TableCell>{user.phone || "-"}</TableCell>
                            <TableCell>{user.city || "-"}</TableCell>
                            <TableCell>{user.specialty || "-"}</TableCell>
                            <TableCell className="max-w-[200px] truncate" title={user.experience}>
                              {user.experience || "-"}
                            </TableCell>
                            <TableCell>
                              {new Date(user.created_at).toLocaleDateString("ru-RU")}
                            </TableCell>
                            <TableCell>
                              <Button size="sm" variant="outline" onClick={(e) => { e.stopPropagation(); openUserModal(user); }}>
                                Подробнее
                              </Button>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>

                {/* Mobile/Tablet cards grid */}
                <div className="lg:hidden grid gap-4 sm:grid-cols-2">
                  {users.map((user, idx) => (
                    <motion.div
                      key={user.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: idx * 0.05 }}
                      onClick={() => openUserModal(user)}
                    >
                      <Card className="cursor-pointer hover:shadow-md transition-shadow">
                        <CardContent className="p-4 space-y-3">
                          <div className="flex items-center gap-3">
                            <div className="h-12 w-12 rounded-full bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-bold">
                              {(user.first_name?.[0] || user.username?.[0] || "U").toUpperCase()}
                            </div>
                            <div>
                              <h3 className="font-medium">{user.first_name} {user.last_name}</h3>
                              <p className="text-sm text-gray-500">@{user.username || "-"}</p>
                            </div>
                          </div>

                          <div className="space-y-1 text-sm">
                            {user.phone && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <Phone className="h-3 w-3" /> {user.phone}
                              </div>
                            )}
                            {user.city && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <MapPin className="h-3 w-3" /> {user.city}
                              </div>
                            )}
                            {user.specialty && (
                              <div className="flex items-center gap-2 text-gray-500">
                                <Briefcase className="h-3 w-3" /> {user.specialty}
                              </div>
                            )}
                          </div>

                          <div className="flex items-center gap-2 text-xs text-gray-400">
                            <Calendar className="h-3 w-3" />
                            {new Date(user.created_at).toLocaleDateString("ru-RU")}
                          </div>
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </main>
      </motion.div>

      {/* User Details Modal */}
      <Dialog open={userModalOpen} onOpenChange={setUserModalOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-3">
              <div className="h-12 w-12 rounded-full bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-bold">
                {(selectedUser?.first_name?.[0] || selectedUser?.username?.[0] || "U").toUpperCase()}
              </div>
              <div>
                <div>{selectedUser?.first_name} {selectedUser?.last_name}</div>
                <div className="text-sm font-normal text-gray-500">@{selectedUser?.username || "-"}</div>
              </div>
            </DialogTitle>
          </DialogHeader>

          {selectedUser && (
            <div className="space-y-4">
              <div className="grid gap-4">
                <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                  <MessageCircle className="h-5 w-5 text-blue-500" />
                  <div>
                    <div className="text-sm text-gray-500">Telegram ID</div>
                    <div className="font-medium">{selectedUser.telegram_id}</div>
                  </div>
                </div>

                {selectedUser.username && (
                  <a
                    href={`https://t.me/${selectedUser.username}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-3 p-3 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors"
                  >
                    <ExternalLink className="h-5 w-5 text-blue-500" />
                    <div>
                      <div className="text-sm text-blue-600">Открыть в Telegram</div>
                      <div className="font-medium text-blue-700">t.me/{selectedUser.username}</div>
                    </div>
                  </a>
                )}

                {selectedUser.phone && (
                  <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                    <Phone className="h-5 w-5 text-green-500" />
                    <div>
                      <div className="text-sm text-gray-500">Телефон</div>
                      <div className="font-medium">{selectedUser.phone}</div>
                    </div>
                  </div>
                )}

                {selectedUser.city && (
                  <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                    <MapPin className="h-5 w-5 text-red-500" />
                    <div>
                      <div className="text-sm text-gray-500">Город</div>
                      <div className="font-medium">{selectedUser.city}</div>
                    </div>
                  </div>
                )}

                {selectedUser.specialty && (
                  <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                    <Briefcase className="h-5 w-5 text-purple-500" />
                    <div>
                      <div className="text-sm text-gray-500">Специальность</div>
                      <div className="font-medium">{selectedUser.specialty}</div>
                    </div>
                  </div>
                )}

                {selectedUser.experience && (
                  <div className="p-3 bg-gray-50 rounded-xl">
                    <div className="text-sm text-gray-500 mb-1">Опыт работы</div>
                    <div className="font-medium">{selectedUser.experience}</div>
                  </div>
                )}

                <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                  <Calendar className="h-5 w-5 text-orange-500" />
                  <div>
                    <div className="text-sm text-gray-500">Дата регистрации</div>
                    <div className="font-medium">
                      {new Date(selectedUser.created_at).toLocaleDateString("ru-RU", {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setUserModalOpen(false)}>
              Закрыть
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Job Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingJob ? "Редактировать вакансию" : "Добавить вакансию"}
            </DialogTitle>
            <DialogDescription>
              Заполните информацию о вакансии
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="title">Название</Label>
              <Input
                id="title"
                value={jobForm.title}
                onChange={(e) => setJobForm({ ...jobForm, title: e.target.value })}
                placeholder="Например: Требуется электрик"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label>Категория</Label>
                <Select
                  value={jobForm.category}
                  onValueChange={(value) =>
                    setJobForm({ ...jobForm, category: value, subcategory: "" })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Выберите категорию" />
                  </SelectTrigger>
                  <SelectContent>
                    {Object.keys(categories).map((cat) => (
                      <SelectItem key={cat} value={cat}>
                        {cat}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label>Подкатегория</Label>
                <Select
                  value={jobForm.subcategory}
                  onValueChange={(value) =>
                    setJobForm({ ...jobForm, subcategory: value })
                  }
                  disabled={!jobForm.category}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Выберите подкатегорию" />
                  </SelectTrigger>
                  <SelectContent>
                    {jobForm.category &&
                      categories[jobForm.category]?.map((sub) => (
                        <SelectItem key={sub} value={sub}>
                          {sub}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label>Город</Label>
                <Select
                  value={jobForm.city}
                  onValueChange={(value) => setJobForm({ ...jobForm, city: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Выберите город" />
                  </SelectTrigger>
                  <SelectContent>
                    {cities.map((city) => (
                      <SelectItem key={city} value={city}>
                        {city}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="salary">Зарплата</Label>
                <Input
                  id="salary"
                  value={jobForm.salary}
                  onChange={(e) => setJobForm({ ...jobForm, salary: e.target.value })}
                  placeholder="Например: 30000-50000 сом"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="phone">Контактный телефон</Label>
                <Input
                  id="phone"
                  value={jobForm.phone}
                  onChange={(e) => setJobForm({ ...jobForm, phone: e.target.value })}
                  placeholder="+996 XXX XXX XXX"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="company">Компания</Label>
                <Input
                  id="company"
                  value={jobForm.company}
                  onChange={(e) => setJobForm({ ...jobForm, company: e.target.value })}
                  placeholder="Название компании (опционально)"
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Описание</Label>
              <Textarea
                id="description"
                value={jobForm.description}
                onChange={(e) =>
                  setJobForm({ ...jobForm, description: e.target.value })
                }
                placeholder="Описание вакансии, требования и т.д."
                rows={4}
              />
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_active"
                checked={jobForm.is_active}
                onChange={(e) =>
                  setJobForm({ ...jobForm, is_active: e.target.checked })
                }
                className="h-4 w-4"
              />
              <Label htmlFor="is_active">Активна (видна пользователям)</Label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              Отмена
            </Button>
            <Button onClick={handleSaveJob}>
              {editingJob ? "Сохранить" : "Создать"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
