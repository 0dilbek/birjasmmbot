package locales

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/models"
)

type Loc struct {
	lang string
}

func New(lang string) *Loc {
	if lang != "uz" {
		lang = "ru"
	}
	return &Loc{lang: lang}
}

// tr holds all UI strings as [RU, UZ] pairs.
var tr = map[string][2]string{
	// Language selection (bilingual — shown before language is known)
	"lang_choice":  {"🌐 Tilni tanlang / Выберите язык:", "🌐 Tilni tanlang / Выберите язык:"},
	"lang_changed_ru": {"✅ Язык изменён на Русский.", "✅ Til ruscha o'zgartirildi."},
	"lang_changed_uz": {"✅ Язык изменён на O'zbek.", "✅ Til o'zbekchaga o'zgartirildi."},

	// General
	"error_start":    {"Ошибка. Попробуйте /start", "Xatolik. /start ni bosing"},
	"no_access":      {"Нет доступа.", "Ruxsat yo'q."},
	"error_generic":  {"Ошибка. Попробуйте позже.", "Xatolik. Keyinroq urinib ko'ring."},

	// Welcome / main menu
	"welcome_new": {
		"👋 Добро пожаловать в <b>Биржу СММ</b>!\n\nНайди исполнителя или зарабатывай на заказах.",
		"👋 <b>BirjaSMM</b>-ga xush kelibsiz!\n\nIjrochi toping yoki buyurtmalar orqali daromad qiling.",
	},
	"client_welcome":   {"👤 <b>%s</b>, добро пожаловать!", "👤 <b>%s</b>, xush kelibsiz!"},
	"exec_welcome":     {"🎬 <b>%s</b>, добро пожаловать!", "🎬 <b>%s</b>, xush kelibsiz!"},
	"account_blocked":  {"🚫 Ваш аккаунт заблокирован.", "🚫 Hisobingiz bloklangan."},

	// Menu buttons (reply keyboard)
	"btn_create_task":     {"➕ Создать задачу", "➕ Vazifa yaratish"},
	"btn_my_tasks":        {"📋 Мои задачи", "📋 Mening vazifalarim"},
	"btn_executors":       {"⭐ Исполнители", "⭐ Ijrochilar"},
	"btn_profile":         {"👤 Профиль", "👤 Profil"},
	"btn_become_verified": {"⭐ Стать проверенным", "⭐ Tasdiqlangan bo'lish"},
	"btn_find_orders":     {"📋 Найти заказы", "📋 Buyurtma topish"},
	"btn_my_responses":    {"📨 Мои отклики", "📨 Mening takliflarim"},
	"btn_in_progress":     {"⚡ В работе", "⚡ Jarayonda"},
	"btn_subscription":    {"💎 Подписка", "💎 Obuna"},
	"btn_change_lang":     {"🌐 Язык", "🌐 Til"},

	// Role
	"role_client": {"👤 Я заказчик", "👤 Men buyurtmachi"},
	"role_exec":   {"🎬 Я исполнитель", "🎬 Men ijrochi"},

	// Registration
	"enter_name":       {"✏️ Введите ваше <b>имя</b>:", "✏️ <b>Ismingiz</b>ni kiriting:"},
	"name_empty":       {"Имя не может быть пустым:", "Ism bo'sh bo'lishi mumkin emas:"},
	"share_phone":      {"📱 Поделитесь вашим <b>номером телефона</b>:", "📱 <b>Telefon raqamingiz</b>ni ulashing:"},
	"phone_btn":        {"📱 Поделиться номером", "📱 Raqamni ulashish"},
	"enter_business":   {"🏢 Введите <b>название бизнеса</b>:", "🏢 <b>Biznes nomi</b>ni kiriting:"},
	"enter_city":       {"🏙 Введите ваш <b>город</b>:", "🏙 <b>Shahringiz</b>ni kiriting:"},
	"reg_done_client":  {"✅ Регистрация завершена! Добро пожаловать!", "✅ Ro'yxatdan o'tish tugadi! Xush kelibsiz!"},
	"choose_category":  {"🎯 Выберите вашу <b>специализацию</b>:", "🎯 <b>Mutaxassislik</b>ingizni tanlang:"},
	"enter_experience": {"📅 Сколько лет опыта? (введите число):", "📅 Necha yil tajribangiz? (raqam kiriting):"},
	"invalid_experience": {"Введите корректное число лет:", "To'g'ri raqam kiriting:"},
	"enter_portfolio":  {"🔗 Введите ссылки на <b>портфолио</b> (через запятую):", "🔗 <b>Portfolio</b> havolalarini kiriting (vergul bilan):"},
	"enter_description": {"📝 Напишите краткое <b>описание</b> о себе:", "📝 O'zingiz haqida qisqacha <b>tavsif</b> yozing:"},
	"reg_done_exec": {
		"✅ Регистрация завершена!\n\nХотите пройти <b>верификацию</b>? Это повысит доверие заказчиков.",
		"✅ Ro'yxatdan o'tish tugadi!\n\n<b>Tasdiq</b>dan o'tmoqchimisiz? Bu buyurtmachilarning ishonchini oshiradi.",
	},

	// Category labels
	"cat_smm":     {"📱 СММ", "📱 SMM"},
	"cat_video":   {"🎥 Мобилограф", "🎥 Mobilograf"},
	"cat_blogger": {"📣 Блогер", "📣 Blogger"},

	// Category buttons
	"btn_cat_smm":     {"📱 СММ", "📱 SMM"},
	"btn_cat_video":   {"🎥 Мобилограф", "🎥 Mobilograf"},
	"btn_cat_blogger": {"📣 Блогер", "📣 Blogger"},

	// Task list
	"no_tasks":          {"У вас пока нет задач.", "Sizda hali vazifalar yo'q."},
	"open_tasks_title":  {"📋 <b>Открытые задачи:</b>", "📋 <b>Ochiq vazifalar:</b>"},
	"no_open_tasks":     {"Открытых задач пока нет.", "Hozircha ochiq vazifalar yo'q."},
	"filter_label":      {"🔍 Фильтр:", "🔍 Filtr:"},
	"filter_cat_menu":   {"🔍 Фильтр по категории:", "🔍 Kategoriya bo'yicha filtr:"},
	"page_label":        {"Страница %d", "%d-sahifa"},
	"task_not_found":    {"Задача не найдена.", "Vazifa topilmadi."},
	"urgent_pro_only":   {"⚡ Срочные задачи доступны только для <b>Pro</b>-подписчиков.", "⚡ Shoshilinch vazifalar faqat <b>Pro</b> obunachilarga mavjud."},
	"choose_action":     {"Выберите действие:", "Amalni tanlang:"},
	"task_in_progress_msg": {"Задача в работе.", "Vazifa jarayonda."},
	"no_active_tasks":   {"У вас нет активных задач в работе.", "Sizda jarayondagi faol vazifalar yo'q."},
	"active_tasks_title": {"⚡ <b>Задачи в работе:</b>", "⚡ <b>Jarayondagi vazifalar:</b>"},
	"no_verified_exec":  {"Проверенных исполнителей пока нет.", "Hozircha tasdiqlangan ijrochilar yo'q."},
	"exec_city_prompt":  {"🏙 Введите <b>город</b> для поиска исполнителей:", "🏙 Ijrochilarni qidirish uchun <b>shahar</b>ni kiriting:"},
	"task_city_prompt":  {"🏙 Введите <b>город</b> для поиска задач:", "🏙 Vazifalarni qidirish uchun <b>shahar</b>ni kiriting:"},
	"exec_city_title":   {"🏙 Показываю исполнителей из города <b>%s</b>:", "🏙 <b>%s</b> shahridan ijrochilar:"},
	"task_city_title":   {"🏙 Задачи из города <b>%s</b>:", "🏙 <b>%s</b> shahridagi vazifalar:"},

	// Inline buttons
	"btn_respond":       {"✍️ Откликнуться", "✍️ Taklif yuborish"},
	"btn_back":          {"⬅️ Назад", "⬅️ Orqaga"},
	"btn_view":          {"👁 Смотреть", "👁 Ko'rish"},
	"btn_view_detail":   {"👁 Подробнее", "👁 Batafsil"},
	"btn_cancel_inline": {"❌ Отменить", "❌ Bekor qilish"},
	"btn_view_profile":  {"👁 Смотреть профиль", "👁 Profilni ko'rish"},
	"btn_select_exec":   {"✅ Выбрать", "✅ Tanlash"},
	"btn_complete_task": {"✅ Завершить задачу", "✅ Vazifani yakunlash"},
	"btn_withdraw_resp": {"↩️ Отозвать отклик", "↩️ Taklifni qaytarish"},
	"btn_cancel_task":   {"❌ Отменить задачу", "❌ Vazifani bekor qilish"},
	"btn_skip":          {"⏭ Пропустить", "⏭ O'tkazib yuborish"},
	"btn_verif_start":   {"🎥 Пройти верификацию", "🎥 Tasdiqdan o'tish"},
	"btn_verif_doc":     {"📄 Отправить документ", "📄 Hujjat yuborish"},
	"btn_cancel":        {"Отмена", "Bekor qilish"},
	"btn_next_page":     {"➡️ Следующая страница", "➡️ Keyingi sahifa"},

	// Filter buttons
	"btn_filter_city":        {"🏙 По городу", "🏙 Shahar bo'yicha"},
	"btn_filter_all":         {"🔄 Все задачи", "🔄 Barcha vazifalar"},
	"btn_filter_all_exec":    {"🔄 Все", "🔄 Hammasi"},
	"btn_filter_budget_low":  {"💰 До 100к", "💰 100 mingacha"},
	"btn_filter_budget_mid":  {"💰 100к-500к", "💰 100к-500к"},
	"btn_filter_budget_high": {"💰 500к+", "💰 500к+"},

	// Task create
	"only_clients":        {"Создавать задачи могут только заказчики.", "Vazifalar faqat buyurtmachilar tomonidan yaratilishi mumkin."},
	"enter_task_title":    {"📝 Введите <b>название задачи</b> (до 50 символов):", "📝 Vazifa <b>nomini</b> kiriting (50 belgigacha):"},
	"title_length":        {"Название должно быть от 1 до 50 символов:", "Nom 1 dan 50 belgigacha bo'lishi kerak:"},
	"enter_task_desc":     {"📋 Введите <b>описание задачи</b> (до 300 символов):", "📋 Vazifa <b>tavsifini</b> kiriting (300 belgigacha):"},
	"desc_too_long":       {"Описание слишком длинное. До 300 символов:", "Tavsif juda uzun. 300 belgigacha:"},
	"choose_task_cat":     {"🎯 Выберите <b>категорию</b>:", "🎯 <b>Kategoriya</b>ni tanlang:"},
	"choose_budget":       {"💰 Выберите <b>тип бюджета</b>:", "💰 <b>Byudjet turini</b> tanlang:"},
	"enter_amount":        {"💵 Введите <b>сумму</b> (в сумах):", "💵 <b>Miqdorni</b> kiriting (so'mda):"},
	"enter_min_amount":    {"💵 Введите <b>минимальную сумму</b>:", "💵 <b>Minimal miqdorni</b> kiriting:"},
	"enter_max_amount":    {"💵 Введите <b>максимальную сумму</b>:", "💵 <b>Maksimal miqdorni</b> kiriting:"},
	"enter_deadline":      {"📅 Введите <b>срок</b> выполнения (ДД.ММ.ГГГГ):", "📅 Bajarish <b>muddatini</b> kiriting (KK.OO.YYYY):"},
	"invalid_date":        {"Неверный формат. Используйте ДД.ММ.ГГГГ (например: 25.05.2025):", "Noto'g'ri format. KK.OO.YYYY (masalan: 25.05.2025):"},
	"enter_refs":          {"🔗 Добавьте <b>референсы</b> (необязательно):", "🔗 <b>Referenslar</b> qo'shing (ixtiyoriy):"},
	"is_urgent":           {"⚡ Задача <b>срочная</b>?", "⚡ Vazifa <b>shoshilinchmi</b>?"},
	"invalid_amount":      {"Введите корректную сумму:", "To'g'ri miqdorni kiriting:"},
	"task_confirm_header": {"📋 <b>Подтверждение задачи</b>", "📋 <b>Vazifani tasdiqlash</b>"},
	"confirm_name":        {"📌 <b>Название:</b>", "📌 <b>Nomi:</b>"},
	"confirm_desc":        {"📝 <b>Описание:</b>", "📝 <b>Tavsif:</b>"},
	"confirm_cat":         {"🎯 <b>Категория:</b>", "🎯 <b>Kategoriya:</b>"},
	"confirm_budget":      {"💰 <b>Бюджет:</b>", "💰 <b>Byudjet:</b>"},
	"confirm_deadline":    {"📅 <b>Срок:</b>", "📅 <b>Muddat:</b>"},
	"confirm_refs":        {"🔗 <b>Референсы:</b>", "🔗 <b>Referenslar:</b>"},
	"confirm_urgent":      {"⚡ <b>Срочно:</b>", "⚡ <b>Shoshilinch:</b>"},
	"btn_publish":         {"✅ Опубликовать", "✅ E'lon qilish"},
	"btn_edit":            {"✏️ Изменить", "✏️ O'zgartirish"},
	"btn_cancel_action":   {"❌ Отмена", "❌ Bekor qilish"},
	"task_published":      {"✅ Задача #%d опубликована!", "✅ Vazifa #%d e'lon qilindi!"},
	"task_error":          {"Ошибка при создании задачи. Попробуйте позже.", "Vazifa yaratishda xatolik. Keyinroq urinib ko'ring."},
	"urgent_yes":          {"⚡ Да", "⚡ Ha"},
	"urgent_no_label":     {"Нет", "Yo'q"},
	"btn_urgent_yes":      {"⚡ Да, срочно", "⚡ Ha, shoshilinch"},
	"btn_urgent_no":       {"📅 Нет", "📅 Yo'q"},
	"btn_budget_fixed":    {"💵 Фикс", "💵 Belgilangan"},
	"btn_budget_range":    {"📊 Диапазон", "📊 Diapazon"},
	"btn_budget_negotiable": {"🤝 Договорной", "🤝 Kelishiladi"},
	"new_task_notify":     {"🔔 Новая задача%s в категории <b>%s</b>!\n\n<b>%s</b>\n%s", "🔔 Yangi vazifa%s\nKategoriya: <b>%s</b>\n\n<b>%s</b>\n%s"},

	// Budget
	"sum_label":          {"сум", "so'm"},
	"budget_negotiable":  {"Договорной", "Kelishiladi"},

	// Task format
	"task_fmt": {
		"📌 <b>%s</b>%s\n\n📝 %s\n🎯 %s\n💰 %s\n📅 Срок: %s\n🔗 %s\n👤 Заказчик: %s\n💬 Откликов: %d/%d",
		"📌 <b>%s</b>%s\n\n📝 %s\n🎯 %s\n💰 %s\n📅 Muddat: %s\n🔗 %s\n👤 Buyurtmachi: %s\n💬 Takliflar: %d/%d",
	},
	"resp_price_line": {"\n💰 Цена: %d сум", "\n💰 Narx: %d so'm"},

	// Response
	"no_responses":        {"Откликов пока нет.", "Hali takliflar yo'q."},
	"only_executors":      {"Откликаться могут только исполнители.", "Taklif faqat ijrochilar yuborishi mumkin."},
	"task_unavailable":    {"Задача недоступна для отклика.", "Vazifa taklif uchun mavjud emas."},
	"resp_limit_exceeded": {"❌ Лимит откликов на эту задачу исчерпан.", "❌ Bu vazifaga takliflar chegarasi to'ldi."},
	"already_responded":   {"Вы уже откликнулись на эту задачу.", "Siz allaqachon bu vazifaga taklif yuborgansiz."},
	"free_limit_exceeded": {"❌ Лимит бесплатных откликов исчерпан. Оформите подписку.", "❌ Bepul takliflar chegarasi tugadi. Obunani rasmiylashtirisng."},
	"write_response":      {"✍️ Напишите ваш <b>отклик</b>:", "✍️ <b>Taklifingiz</b>ni yozing:"},
	"free_resps_hint":     {"\n\n💡 Осталось бесплатных откликов: <b>%d</b>", "\n\n💡 Qolgan bepul takliflar: <b>%d</b>"},
	"msg_empty":           {"Сообщение не может быть пустым:", "Xabar bo'sh bo'lishi mumkin emas:"},
	"enter_price":         {"💰 Укажите вашу <b>цену</b> (необязательно):", "💰 <b>Narxingiz</b>ni kiriting (ixtiyoriy):"},
	"invalid_price":       {"Введите корректную сумму или нажмите «Пропустить»:", "To'g'ri miqdorni kiriting yoki «O'tkazib yuborish»ni bosing:"},
	"resp_error":          {"Ошибка при отправке отклика. Попробуйте позже.", "Taklif yuborishda xatolik. Keyinroq urinib ko'ring."},
	"resp_sent":           {"✅ Отклик отправлен!", "✅ Taklif yuborildi!"},
	"new_resp_notify":     {"📬 Новый отклик на задачу <b>%s</b> от <b>%s</b>", "📬 <b>%s</b> vazifasiga <b>%s</b> tomonidan yangi taklif"},
	"no_my_responses":     {"У вас пока нет откликов.", "Sizda hali takliflar yo'q."},
	"resp_not_found":      {"Отклик не найден.", "Taklif topilmadi."},
	"resp_withdrawn_msg":  {"↩️ Отклик отозван.", "↩️ Taklif qaytarildi."},
	"resp_cant_withdraw":  {"Этот отклик уже нельзя отозвать.", "Bu taklifni qaytarib bo'lmaydi."},

	// Resp status labels
	"resp_waiting":  {"⏳ Ожидает", "⏳ Kutilmoqda"},
	"resp_accepted": {"✅ Принят", "✅ Qabul qilindi"},
	"resp_rejected": {"❌ Отклонён", "❌ Rad etildi"},

	// Review / complete
	"task_complete_review": {"✅ Задача завершена!\n\nОцените работу исполнителя <b>%s</b>:", "✅ Vazifa yakunlandi!\n\n<b>%s</b> ijrochisining ishini baholang:"},
	"task_complete_simple": {"✅ Задача завершена!", "✅ Vazifa yakunlandi!"},
	"write_review":         {"💬 Напишите <b>отзыв</b> (необязательно):", "💬 <b>Fikr</b> yozing (ixtiyoriy):"},
	"thanks_review":        {"⭐ Спасибо за отзыв!", "⭐ Fikr uchun rahmat!"},
	"exec_selected":        {"✅ Исполнитель <b>%s</b> выбран! Задача переведена в статус «В работе».", "✅ Ijrochi <b>%s</b> tanlandi! Vazifa «Jarayonda» holatiga o'tkazildi."},
	"exec_not_found":       {"Отклик не найден.", "Taklif topilmadi."},
	"rejected_notify":      {"ℹ️ По задаче <b>%s</b> выбран другой исполнитель. Удачи в следующий раз!", "ℹ️ <b>%s</b> vazifasi uchun boshqa ijrochi tanlandi. Keyingi safar omad!"},
	"selected_notify":      {"🎉 Вы выбраны исполнителем по задаче <b>%s</b>!\n\nКонтакт заказчика:\n👤 %s\nTelegram: @%s\n📱 %s", "🎉 Siz <b>%s</b> vazifasi uchun ijrochi etib tanlandingiz!\n\nBuyurtmachi kontakti:\n👤 %s\nTelegram: @%s\n📱 %s"},

	// Subscription
	"sub_menu": {
		"💎 <b>Подписки</b>\n\n<b>Basic</b> — доступ ко всем заказам\n<b>Pro</b> — приоритет в откликах + ⚡ срочные задачи\n\nВыберите тариф:",
		"💎 <b>Obunalar</b>\n\n<b>Basic</b> — barcha buyurtmalarga kirish\n<b>Pro</b> — takliflarda ustuvorlik + ⚡ shoshilinch vazifalar\n\nTarifni tanlang:",
	},
	"sub_only_exec": {"Подписки доступны только для исполнителей.", "Obunalar faqat ijrochilar uchun."},
	"btn_sub_basic": {"💡 Basic (30 дней)", "💡 Basic (30 kun)"},
	"btn_sub_pro":   {"💎 Pro (30 дней)", "💎 Pro (30 kun)"},

	// Payment
	"payment_prompt": {
		"💳 <b>Оплата подписки %s</b>\n\nСумма: <b>%d сум</b>\n\nПереведите на карту:\n<code>%s</code>\n\nПосле оплаты отправьте <b>скриншот чека</b> (фото, файл или текст сообщением):",
		"💳 <b>%s obuna to'lovi</b>\n\nMiqdor: <b>%d so'm</b>\n\nKartaga o'tkazing:\n<code>%s</code>\n\nTo'lovdan so'ng <b>kvitansiya skrinshotini</b> yuboring (foto, fayl yoki matn):",
	},
	"payment_send_receipt": {"Пожалуйста, отправьте фото, файл или текстовый чек:", "Iltimos, foto, fayl yoki matn kvitansiyasi yuboring:"},
	"payment_error":        {"Ошибка при сохранении чека. Попробуйте позже.", "Kvitansiyani saqlashda xatolik. Keyinroq urinib ko'ring."},
	"payment_saved":        {"✅ Чек получен! Ожидайте подтверждения администратора.", "✅ Kvitansiya qabul qilindi! Administrator tasdig'ini kuting."},
	"payment_approved": {
		"✅ Ваш платёж подтверждён! Подписка <b>%s</b> активирована на 30 дней.",
		"✅ To'lovingiz tasdiqlandi! <b>%s</b> obuna 30 kunga faollashtirildi.",
	},
	"payment_rejected": {
		"❌ Ваш платёж отклонён. Обратитесь к администратору за уточнением.",
		"❌ To'lovingiz rad etildi. Aniqlik uchun administratorga murojaat qiling.",
	},

	// Verification
	"verif_code_prompt": {
		"🎥 Запишите видео, покажите лицо и назовите код: <b>%s</b>\n\nЗатем отправьте видео в этот чат.",
		"🎥 Video yozing, yuzingizni ko'rsating va kodni ayting: <b>%s</b>\n\nKeyin videoni shu chatga yuboring.",
	},
	"verif_send_video":        {"Пожалуйста, отправьте видео.", "Iltimos, video yuboring."},
	"verif_sent":              {"✅ Видео отправлено на проверку! Мы уведомим вас о результате.", "✅ Video tekshiruv uchun yuborildi! Natija haqida xabardor qilamiz."},
	"verif_approved_exec":     {"✅ Ваша верификация одобрена! Теперь вы верифицированный исполнитель.", "✅ Tasdiqingiz tasdiqlandi! Endi siz tasdiqlangan ijrochisiz."},
	"verif_approved_client":   {"✅ Ваш бизнес-аккаунт верифицирован!", "✅ Biznes hisobingiz tasdiqlandi!"},
	"verif_rejected_user":     {"❌ Ваша верификация отклонена. Попробуйте снова.", "❌ Tasdiqingiz rad etildi. Qaytadan urinib ko'ring."},
	"already_verified":        {"✅ Вы уже верифицированы.", "✅ Siz allaqachon tasdiqlangansissiz."},
	"verif_pending_status":    {"⏳ Ваша заявка уже на рассмотрении.", "⏳ So'rovingiz ko'rib chiqilmoqda."},
	"verif_client_offer": {
		"📄 Для верификации отправьте документ, подтверждающий ваш бизнес (фото, скан, PDF).\n\nЭто повысит доверие исполнителей к вашим задачам.",
		"📄 Tasdiqlash uchun biznesingizni tasdiqlovchi hujjat yuboring (foto, skan, PDF).\n\nBu ijrochilarning vazifalaringizga ishonchini oshiradi.",
	},
	"verif_client_doc_prompt": {"📎 Отправьте документ (фото или файл):", "📎 Hujjat yuboring (foto yoki fayl):"},
	"verif_client_send_doc":   {"Пожалуйста, отправьте фото или файл документа.", "Iltimos, foto yoki hujjat faylini yuboring."},
	"verif_client_sent":       {"✅ Документ отправлен на проверку! Ожидайте подтверждения.", "✅ Hujjat tekshiruv uchun yuborildi! Tasdig'ini kuting."},
	"verif_client_approved": {
		"✅ Ваш бизнес-аккаунт верифицирован! Теперь вы отображаетесь как проверенный заказчик.",
		"✅ Biznes hisobingiz tasdiqlandi! Endi siz tasdiqlangan buyurtmachi sifatida ko'rinasiz.",
	},
	"verif_client_rejected": {
		"❌ Верификация отклонена. Попробуйте отправить другой документ.",
		"❌ Tasdiq rad etildi. Boshqa hujjat yuborib ko'ring.",
	},

	// Profile
	"profile_not_found":    {"Профиль не найден.", "Profil topilmadi."},
	"profile_verified":     {"✅ Да", "✅ Ha"},
	"profile_not_verified": {"Нет", "Yo'q"},
	"client_profile_text": {
		"👤 <b>Профиль заказчика</b>\n\nИмя: %s\nБизнес: %s\nГород: %s\nПроверенный: %s\nТелефон: %s",
		"👤 <b>Buyurtmachi profili</b>\n\nIsm: %s\nBiznes: %s\nShahar: %s\nTasdiqlangan: %s\nTelefon: %s",
	},
	"exec_profile_text": {
		"🎬 <b>Профиль исполнителя</b>\n\nИмя: %s\nСпециализация: %s\nГород: %s\nОпыт: %d лет\nРейтинг: ⭐ %.1f (%d заказов)\n%s\n\nПодписка: %s\nБесплатных откликов: %d\n\nПортфолио: %s",
		"🎬 <b>Ijrochi profili</b>\n\nIsm: %s\nMutaxassislik: %s\nShahar: %s\nTajriba: %d yil\nReyting: ⭐ %.1f (%d buyurtma)\n%s\n\nObuna: %s\nBepul takliflar: %d\n\nPortfolio: %s",
	},
	"exec_verified_label": {"✅ Верифицирован", "✅ Tasdiqlangan"},
	"exec_not_verified_label": {"❌ Не верифицирован", "❌ Tasdiqlanmagan"},
	"sub_active_label":    {"%s до %s", "%s — %s gacha"},
	"sub_none_label":      {"Нет", "Yo'q"},

	// Task cancel
	"task_cancelled_msg":    {"❌ Задача отменена.", "❌ Vazifa bekor qilindi."},
	"task_cancel_forbidden": {"Задачу нельзя отменить — она уже не открыта.", "Vazifani bekor qilib bo'lmaydi — u ochiq emas."},

	// Active task view
	"active_task_text": {
		"📋 <b>%s</b>\n%s\n\n👤 Заказчик: %s%s\n\n📅 Назначено: %s",
		"📋 <b>%s</b>\n%s\n\n👤 Buyurtmachi: %s%s\n\n📅 Tayinlangan: %s",
	},
	"client_tg_line":    {"\n📱 Telegram: @%s", "\n📱 Telegram: @%s"},
	"client_phone_line": {"\n☎️ Телефон: %s", "\n☎️ Telefon: %s"},
}

// menuActions maps button texts (both languages) to action keys.
var menuActions = map[string]string{
	// Russian
	"➕ Создать задачу":     "create_task",
	"📋 Мои задачи":         "my_tasks",
	"⭐ Исполнители":        "executors",
	"👤 Профиль":            "profile",
	"⭐ Стать проверенным":  "become_verified",
	"📋 Найти заказы":       "find_orders",
	"📨 Мои отклики":        "my_responses",
	"⚡ В работе":           "in_progress",
	"💎 Подписка":           "subscription",
	"🌐 Язык":               "change_lang",
	// Uzbek
	"➕ Vazifa yaratish":     "create_task",
	"📋 Mening vazifalarim":  "my_tasks",
	"⭐ Ijrochilar":          "executors",
	"👤 Profil":              "profile",
	"⭐ Tasdiqlangan bo'lish": "become_verified",
	"📋 Buyurtma topish":     "find_orders",
	"📨 Mening takliflarim":  "my_responses",
	"⚡ Jarayonda":           "in_progress",
	"💎 Obuna":               "subscription",
	"🌐 Til":                 "change_lang",
}

// MenuAction returns the action key for a reply keyboard button text (language-agnostic).
func MenuAction(text string) string {
	return menuActions[text]
}

// T translates key to the Loc's language. Supports fmt.Sprintf args.
func (l *Loc) T(key string, args ...interface{}) string {
	if pair, ok := tr[key]; ok {
		s := pair[0] // default RU
		if l.lang == "uz" {
			s = pair[1]
		}
		if len(args) > 0 {
			return fmt.Sprintf(s, args...)
		}
		return s
	}
	return key
}

// ─── Keyboard methods ──────────────────────────────────────────────────────────

func LangChoice() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇺🇿 O'zbek", "lang:uz"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
		),
	)
}

func (l *Loc) RoleChoice() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("role_client"), "role:client"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("role_exec"), "role:executor"),
		),
	)
}

func (l *Loc) RequestPhone() tgbotapi.ReplyKeyboardMarkup {
	btn := tgbotapi.NewKeyboardButtonContact(l.T("phone_btn"))
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(btn))
}

func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}

func (l *Loc) ExecutorCategory() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_smm"), "exec_cat:smm"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_video"), "exec_cat:video"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_blogger"), "exec_cat:blogger"),
		),
	)
}

func (l *Loc) TaskCategory() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_smm"), "task_cat:smm"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_video"), "task_cat:video"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_blogger"), "task_cat:blogger"),
		),
	)
}

func (l *Loc) BudgetType() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_budget_fixed"), "budget:fixed"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_budget_range"), "budget:range"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_budget_negotiable"), "budget:negotiable"),
		),
	)
}

func (l *Loc) UrgentChoice() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_urgent_yes"), "urgent:yes"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_urgent_no"), "urgent:no"),
		),
	)
}

func (l *Loc) TaskConfirm() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_publish"), "task_confirm:yes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_edit"), "task_confirm:edit"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cancel_action"), "task_confirm:cancel"),
		),
	)
}

func (l *Loc) SkipRefs() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_skip"), "skip:refs"),
		),
	)
}

func (l *Loc) SkipPrice() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_skip"), "skip:price"),
		),
	)
}

func (l *Loc) SkipReviewComment() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_skip"), "skip:comment"),
		),
	)
}

func (l *Loc) ClientMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_create_task")),
			tgbotapi.NewKeyboardButton(l.T("btn_my_tasks")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_executors")),
			tgbotapi.NewKeyboardButton(l.T("btn_profile")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_become_verified")),
			tgbotapi.NewKeyboardButton(l.T("btn_change_lang")),
		),
	)
}

func (l *Loc) ExecutorMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_find_orders")),
			tgbotapi.NewKeyboardButton(l.T("btn_my_responses")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_in_progress")),
			tgbotapi.NewKeyboardButton(l.T("btn_profile")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(l.T("btn_subscription")),
			tgbotapi.NewKeyboardButton(l.T("btn_change_lang")),
		),
	)
}

func (l *Loc) VerifOffer() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_verif_start"), "verif:start"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_skip"), "verif:skip"),
		),
	)
}

func (l *Loc) VerifVideo() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_skip"), "verif:skip"),
		),
	)
}

func (l *Loc) ClientVerifOffer() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_verif_doc"), "client_verif:start"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cancel"), "client_verif:cancel"),
		),
	)
}

func (l *Loc) TaskActions(taskID int64, canRespond bool) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	if canRespond {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_respond"), fmt.Sprintf("respond:%d", taskID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(l.T("btn_back"), "back:tasks"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (l *Loc) ResponseActions(respID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_select_exec"), fmt.Sprintf("select_exec:%d", respID)),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_back"), "back:responses"),
		),
	)
}

func (l *Loc) CompleteTask(taskID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_complete_task"), fmt.Sprintf("complete:%d", taskID)),
		),
	)
}

func (l *Loc) WithdrawResponse(respID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_withdraw_resp"), fmt.Sprintf("withdraw_resp:%d", respID)),
		),
	)
}

func (l *Loc) CancelTask(taskID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cancel_task"), fmt.Sprintf("cancel_task:%d", taskID)),
		),
	)
}

func (l *Loc) SubscriptionMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_sub_basic"), "sub:basic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_sub_pro"), "sub:pro"),
		),
	)
}

func (l *Loc) FilterMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_smm"), "filter_cat:smm"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_video"), "filter_cat:video"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_blogger"), "filter_cat:blogger"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_budget_low"), "filter_budget:100000"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_budget_mid"), "filter_budget:100000-500000"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_budget_high"), "filter_budget:500000"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_city"), "filter_city:task"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_all"), "filter_cat:all"),
		),
	)
}

func (l *Loc) ExecutorFilterMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_smm"), "exec_filter:smm"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_video"), "exec_filter:video"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_cat_blogger"), "exec_filter:blogger"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_city"), "exec_filter:city"),
			tgbotapi.NewInlineKeyboardButtonData(l.T("btn_filter_all_exec"), "exec_filter:all"),
		),
	)
}

func RatingKeyboard(taskID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐1", fmt.Sprintf("rate:1:%d", taskID)),
			tgbotapi.NewInlineKeyboardButtonData("⭐2", fmt.Sprintf("rate:2:%d", taskID)),
			tgbotapi.NewInlineKeyboardButtonData("⭐3", fmt.Sprintf("rate:3:%d", taskID)),
			tgbotapi.NewInlineKeyboardButtonData("⭐4", fmt.Sprintf("rate:4:%d", taskID)),
			tgbotapi.NewInlineKeyboardButtonData("⭐5", fmt.Sprintf("rate:5:%d", taskID)),
		),
	)
}

func (l *Loc) TaskNextPage(page int, cat string, budgetMin int64) tgbotapi.InlineKeyboardMarkup {
	catParam := cat
	if catParam == "" {
		catParam = "_"
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				l.T("btn_next_page"),
				fmt.Sprintf("task_page:%d:%s:%d", page, catParam, budgetMin),
			),
		),
	)
}

func (l *Loc) ExecNextPage(page int, cat string) tgbotapi.InlineKeyboardMarkup {
	catParam := cat
	if catParam == "" {
		catParam = "_"
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				l.T("btn_next_page"),
				fmt.Sprintf("exec_page:%d:%s", page, catParam),
			),
		),
	)
}

// ─── Format helpers ───────────────────────────────────────────────────────────

func (l *Loc) BudgetText(t *models.Task) string {
	sum := l.T("sum_label")
	switch t.BudgetType {
	case models.BudgetFixed:
		if t.BudgetFrom != nil {
			return fmt.Sprintf("%d %s", *t.BudgetFrom, sum)
		}
	case models.BudgetRange:
		if t.BudgetFrom != nil && t.BudgetTo != nil {
			return fmt.Sprintf("%d – %d %s", *t.BudgetFrom, *t.BudgetTo, sum)
		}
	case models.BudgetNegotiable:
		return l.T("budget_negotiable")
	}
	return "—"
}

func (l *Loc) CatLabel(cat string) string {
	switch cat {
	case "smm":
		return l.T("cat_smm")
	case "video":
		return l.T("cat_video")
	case "blogger":
		return l.T("cat_blogger")
	}
	return cat
}

func (l *Loc) RespLabel(s models.RespStatus) string {
	switch s {
	case models.RespActive:
		return l.T("resp_waiting")
	case models.RespAccepted:
		return l.T("resp_accepted")
	case models.RespRejected:
		return l.T("resp_rejected")
	}
	return string(s)
}

func (l *Loc) FormatTask(t *models.Task) string {
	urgent := ""
	if t.IsUrgent {
		urgent = " ⚡"
	}
	return fmt.Sprintf(
		l.T("task_fmt"),
		t.Title, urgent,
		t.Description,
		l.CatLabel(string(t.Category)),
		l.BudgetText(t),
		t.Deadline.Format("02.01.2006"),
		ifEmpty(t.Refs, "—"),
		t.ClientName,
		t.ResponseCount, t.MaxResponses,
	)
}

func (l *Loc) FormatResponse(r *models.Response) string {
	badges := ""
	if r.IsPro {
		badges += " 💎"
	}
	if r.IsVerified {
		badges += " ✅"
	}
	price := ""
	if r.ProposedPrice != nil {
		price = l.T("resp_price_line", *r.ProposedPrice)
	}
	return fmt.Sprintf("👤 <b>%s</b>%s\n⭐ %.1f\n\n%s%s",
		r.ExecutorName, badges, r.Rating, r.Message, price)
}

func ifEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
