// The English message catalogue and the source of truth for the key set: every
// other locale is a Partial of these keys (es.ts), and `t()` falls back here for
// any key a translation has not filled yet. Flat dotted keys so `MessageKey`
// gives autocomplete and exhaustiveness for free - no nested-path type magic.
// `{name}`-style placeholders are filled by t(key, params). Grouped by area.

export const en = {
	// Language switcher + locale names.
	'lang.en': 'English',
	'lang.es': 'Español',
	'lang.switch': 'Language',

	// Top nav + global chrome.
	'nav.races': 'Races',
	'nav.sell': 'Sell',
	'nav.myListings': 'My listings',
	'nav.inbox': 'Inbox',
	// Header unread badge (pluralized via Intl.PluralRules, see messages.ts).
	'nav.inboxUnread.one': '{n} unread message',
	'nav.inboxUnread.other': '{n} unread messages',
	'nav.login': 'Log in',
	'nav.register': 'Register',
	'nav.logout': 'Log out',
	'nav.skipToContent': 'Skip to content',
	'banner.verifyEmail': 'Verify your email to unlock selling and chat.',
	'banner.resend': 'Resend email',
	'banner.verifySent': 'Verification email sent - check your inbox (and spam).',
	// Spanish-suggestion banner. Rendered with the *suggested* locale's translator
	// (es), so it shows English now and Spanish once es.ts is filled in M8.2.
	'banner.suggestText': 'This page is available in Spanish.',
	'banner.suggestAccept': 'View in Spanish',
	'banner.suggestDismiss': 'Dismiss',
	'footer.tagline': 'Zero commission, EU-wide. Non-profit by design.',
	'footer.github': 'GitHub',
	'footer.terms': 'Terms',
	'footer.privacy': 'Privacy',
	'footer.contact': 'Contact',

	// Home.
	'home.title': 'Bibseller - race bibs find new runners',
	'home.metaDescription':
		"Non-profit, EU-wide marketplace connecting runners who can't start with runners who missed registration.",
	'home.heroTitle': 'Race bibs find',
	'home.heroTitleHighlight': 'new runners',
	'home.tagline':
		"Injured? Plans changed? Missed registration? A zero-commission marketplace that connects sellers and buyers of race bibs - always within each race's own rules.",
	'home.searchPlaceholder': 'Search a race or city…',
	'home.search': 'Search',
	'home.browseAll': 'or browse all races',
	'home.moreFilters': 'More filters',
	'home.apiUnreachable': 'API unreachable - run',
	'home.upcoming': 'Upcoming races',
	'home.browseAllRaces': 'Browse all races',
	'home.howTitle': 'How it works',
	'home.step1Title': 'List your bib',
	'home.step1Desc': 'Post your race entry in minutes. The transfer policy shows up front.',
	'home.step2Title': 'Message the seller',
	'home.step2Desc': 'Chat securely on Bibseller to agree the details. Zero commission, always.',
	'home.step3Title': 'Arrange the transfer',
	'home.step3Desc': 'Re-register through the organiser where required, or hand over directly.',
	'home.step4Title': 'Confirm the handover',
	'home.step4Desc': 'Both of you confirm in chat. The bib is theirs, and you are done.',
	'home.howNote': 'Some races are chat-only or official resale, set by the race organiser.',
	'home.journeyTitle': 'The buyer and seller journey',
	'home.journeyLead': 'Six simple steps, from listing a bib to the start line.',
	'home.journeySeller': 'Seller',
	'home.journeyBuyer': 'Buyer',
	'home.j1Title': 'Lists the bib',
	'home.j2Title': 'Finds it and messages',
	'home.j3Title': 'Replies and agrees',
	'home.j4Title': 'Sorts the transfer',
	'home.j5Title': 'Hands it over',
	'home.j6Title': 'Toes the start line',
	'home.contactTitle': 'Get in touch',
	'home.contactLead': 'A race to add, a question, or feedback? We would love to hear from you.',
	'home.contactCta': 'Contact us',
	'home.modePlatformSaleName': 'Platform sale',
	'home.modePlatformSaleDesc':
		'The race allows resale: list, chat, and pay securely through the platform.',
	'home.modeOfficialName': 'Official process',
	'home.modeOfficialDesc':
		'The race runs its own name change: we connect you and link the official procedure.',
	'home.modeConnectName': 'Connect only',
	'home.modeConnectDesc':
		'Restricted or unverified races: we provide the chat, the rest stays between you two.',
	'home.underConstruction': 'Under construction - follow the',
	'home.roadmap': 'roadmap',

	// Policy words (moved from $lib/policy.ts - the facts stay there, words live here).
	'policy.label.platform_sale': 'Resale allowed',
	'policy.label.official_only': 'Official transfer',
	'policy.label.connect_only': 'Chat only',
	'policy.label.unknown': 'Policy unverified',
	'policy.disclaimer.platform_sale.title': 'This race allows bib resale.',
	'policy.disclaimer.platform_sale.body':
		'Agree with the seller in chat, then pay securely through the platform - funds are held until the transfer is confirmed. Zero commission.',
	'policy.disclaimer.official_only.title': 'This race runs its own official name-change process.',
	'policy.disclaimer.official_only.body':
		'Find each other and agree on the details here - the transfer itself (and any official fee) goes through the race organizer. The platform never handles money for this race.',
	'policy.disclaimer.connect_only.title': 'This race restricts bib transfers.',
	'policy.disclaimer.connect_only.body':
		"The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything.",
	'policy.disclaimer.unknown.title':
		'Transfer policy not verified yet - treat this race as chat-only.',
	'policy.disclaimer.unknown.body':
		"The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything.",
	'policy.officialLink': 'Official transfer process',

	// Sport filter values (the option value stays the enum; only the label is translated).
	'sport.running': 'Running',
	'sport.trail': 'Trail',
	'sport.triathlon': 'Triathlon',
	'sport.cycling': 'Cycling',
	'sport.obstacle': 'Obstacle',
	'sport.other': 'Other',

	// Races browse + filters.
	'races.title': 'Browse races - Bibseller',
	'races.metaDescription': 'Find race bibs for sale across EU running events.',
	'races.heading': 'Browse races',
	'races.filter.search': 'Search',
	'races.filter.searchPlaceholder': 'Race or city…',
	'races.filter.country': 'Country',
	'races.filter.sport': 'Sport',
	'races.filter.policy': 'Transfer policy',
	'races.filter.distance': 'Distance',
	'races.filter.dateFrom': 'From',
	'races.filter.dateTo': 'Until',
	'races.filtersSummary': 'Filters',
	'races.resultCount.one': '{n} race',
	'races.resultCount.other': '{n} races',
	'races.resultCountMore': 'Showing the first {n} races',
	'races.removeFilter': 'Remove filter: {name}',
	'races.filter.all': 'All',
	'races.filter.submit': 'Filter',
	'races.empty': 'No races match those filters.',
	'races.clearFilters': 'Clear filters',
	'races.nextPage': 'Next page ->',
	'races.mapSummary': 'Map',
	'races.mapHeading': 'Races across Europe',
	'races.mapHint': 'Highlighted countries have races - tap one to filter.',
	'races.mapCountry.one': '{country}: {n} race',
	'races.mapCountry.other': '{country}: {n} races',
	'races.mapCity.one': '{city}: {n} race',
	'races.mapCity.other': '{city}: {n} races',
	'races.mapBack': '<- All of Europe',

	// Race + listing cards. CLDR plural forms via Intl.PluralRules (no ICU lib, D14);
	// en/es use one|other - add .few/.many keys if a locale needs them.
	'raceCard.bibs.one': '{n} bib listed',
	'raceCard.bibs.other': '{n} bibs listed',
	'listingCard.priceOnRequest': 'Price on request',
	'listingCard.belowFace': 'below face value',
	'listingCard.listedBy': 'Listed by {name}',

	// Race detail.
	'raceDetail.title': '{name} - bibs for sale - Bibseller',
	'raceDetail.metaDescription': 'Bibs for {name} ({date}, {city}).',
	'raceDetail.back': 'Back to all races',
	'raceDetail.website': 'Race website',
	'raceDetail.bibsForSale.one': '{n} bib for sale',
	'raceDetail.bibsForSale.other': '{n} bibs for sale',
	'raceDetail.sellCta': 'Sell your bib',
	'raceDetail.empty': 'No bibs listed for this race yet.',
	'raceDetail.emptyHint': 'Selling yours?',

	// Listing detail.
	'listingDetail.title': 'Bib for {name} - Bibseller',
	'listingDetail.back': 'Back to {name}',
	'listingDetail.heading': 'Bib for {name}',
	'listingDetail.unavailable': 'This listing is no longer available ({status}).',
	'listingDetail.listedByOn': 'Listed by {name} on {date}',
	'listingDetail.contact': 'Contact the seller',
	'listingDetail.toMessageSeller': 'to message the seller.',
	'listingDetail.verifyToMessage': 'Verify your email to message the seller.',
	'listingDetail.accountSettings': 'Account settings',
	'listingDetail.ownPre': 'This is your listing - manage it from',
	'listingDetail.yourListings': 'your listings',
	'listingDetail.messageAria': 'Message to the seller',
	'listingDetail.messagePlaceholder': 'Hi - is this bib still available?',
	'listingDetail.ackText':
		"I understand the platform handles no money and takes no responsibility for this transfer - the race's own rules apply.",
	'listingDetail.send': 'Send message',

	// Report a listing / content.
	'report.summary': 'Report this listing',
	'report.reason.forbidden_transfer': 'Forbidden transfer',
	'report.reason.scam': 'Scam',
	'report.reason.offensive': 'Offensive',
	'report.reason.other': 'Other',
	'report.reasonAria': 'Reason for report',
	'report.detailsAria': 'Report details (optional)',
	'report.detailsPlaceholder': 'Details (optional)',
	'report.success': 'Thanks - this listing has been reported.',
	'report.failed': 'Could not file the report. Try again.',
	'report.networkError': 'Network error - try again.',
	'report.submitting': 'Reporting...',
	'report.submit': 'Submit report',

	// Listing CTA (the buy path is an honest disabled stub until M6).
	'listingCta.buy': 'Buy securely - coming soon',
	'listingCta.buyTitle': 'Secure checkout arrives with payments (M6)',

	// Inbox (thread list).
	'inbox.title': 'Inbox - Bibseller',
	'inbox.heading': 'Inbox',
	'inbox.emptyPre': 'No conversations yet. Browse',
	'inbox.emptyRacesLink': 'races',
	'inbox.emptyPost': 'and contact a seller to start one.',
	'role.seller': 'seller',
	'role.buyer': 'buyer',

	// Chat thread.
	'chat.title': 'Chat with {name} - Bibseller',
	'chat.back': 'Back to inbox',
	'chat.about': 'about',
	'chat.block': 'Block',
	'chat.unblock': 'Unblock',
	'chat.sharedImage': 'Shared image',
	'chat.reportMsg': 'report',
	'chat.messageAria': 'Your message',
	'chat.messagePlaceholder': 'Write a message, or attach an image...',
	'chat.attachAria': 'Attach an image (JPEG or PNG)',
	'chat.send': 'Send',
	'chat.sending': 'Sending...',
	'chat.imageTooLarge': 'That image is too large (5 MB max).',
	'chat.tooFast': 'You are sending messages too fast - wait a moment.',
	'chat.sendFailed': 'Could not send your message. Try again.',
	'chat.networkError': 'Network error - check your connection.',
	'chat.blockConfirm': 'Block {name}? Neither of you will be able to message the other.',
	'chat.blocked': 'User blocked.',
	'chat.blockFailed': 'Could not block the user.',
	'chat.unblocked': 'User unblocked.',
	'chat.unblockFailed': 'Could not unblock the user.',
	'chat.networkRetry': 'Network error - try again.',
	'chat.reportConfirm': 'Report this message to the moderators?',
	'chat.messageReported': 'Message reported.',
	'chat.messageReportFailed': 'Could not report the message.',
	'chat.safetySummary': 'Safety options',
	'chat.policyReminder':
		"Reminder: this race restricts transfers - the platform only connects you, so follow the race's own rules and never send money here.",
	'chat.connectionLost': 'Connection lost - retrying...',
	'chat.invalidImage': 'That file could not be read as an image.',
	'chat.unsupportedImage': 'Only JPEG and PNG images are allowed.',
	'chat.blockedSend': 'You cannot send messages in this conversation.',
	'chat.logAria': 'Conversation messages',
	'chat.previewAlt': 'Preview of the attached image',
	'chat.clearImage': 'Remove the attached image',

	// Sell - race search.
	'sell.title': 'Sell a bib - Bibseller',
	'sell.heading': 'Sell a bib',
	'sell.lede': 'Find your race, then list your bib. You set the price (capped at face value).',
	'sell.verifyNotice':
		'Verify your email to publish a listing - you can still find your race below.',
	'sell.searchAria': 'Search races by name or city',
	'sell.emptyPre': 'No upcoming races match. Try another search, or',
	'sell.browseAllLink': 'browse all races',
	'sell.sellHere': 'Sell here',

	// Sell - listing form.
	'sellForm.title': 'List your bib for {name} - Bibseller',
	'sellForm.back': 'Back to race search',
	'sellForm.heading': 'List your bib',
	'sellForm.verifyNotice': 'Verify your email to publish a listing.',
	'sellForm.publish': 'Publish listing',

	// Shared listing form fields (create + edit).
	'listingFields.price': 'Asking price (EUR)',
	'listingFields.pricePlaceholder': 'e.g. 45',
	'listingFields.original': 'Original price / face value (EUR)',
	'listingFields.optional': 'optional',
	'listingFields.hint': 'Enter the face value and your asking price is capped at it - no scalping.',
	'listingFields.description': 'Description',
	'listingFields.descriptionPlaceholder': 'optional - size, pickup details, etc.',

	// My listings.
	'myListings.title': 'My listings - Bibseller',
	'myListings.heading': 'My listings',
	'myListings.emptyPre': 'You have no listings yet.',
	'myListings.listABib': 'List a bib',
	'myListings.edit': 'Edit',
	'myListings.cancel': 'Cancel',
	'myListings.view': 'View',
	'myListings.cancelConfirm': 'Cancel this listing? Buyers will no longer see it.',
	'myListings.cancelled': 'Listing cancelled.',
	'myListings.created': 'Your listing is live.',

	// Listing status labels - the full listings_status_check set. Shared by the
	// my-listings pills and the listing page's unavailable banner.
	'listingStatus.active': 'Active',
	'listingStatus.reserved': 'Reserved',
	'listingStatus.sold': 'Sold',
	'listingStatus.cancelled': 'Cancelled',
	'listingStatus.expired': 'Expired',
	'listingStatus.removed': 'Removed',

	// Edit listing.
	'editListing.title': 'Edit listing - Bibseller',
	'editListing.back': 'Back to my listings',
	'editListing.heading': 'Edit listing',
	'editListing.save': 'Save changes',

	// Auth - shared field labels.
	'auth.email': 'Email',
	'auth.password': 'Password',

	// Log in.
	'login.title': 'Log in - Bibseller',
	'login.forgot': 'Forgot your password?',
	'login.newHere': 'New here?',
	'login.createAccount': 'Create an account',

	// Register.
	'register.title': 'Create account - Bibseller',
	'register.heading': 'Create account',
	'register.displayName': 'Display name',
	'register.haveAccount': 'Already have an account?',

	// Forgot password (request a reset link).
	'forgot.title': 'Reset password - Bibseller',
	'forgot.heading': 'Reset your password',
	'forgot.lede': "Enter your email and we'll send you a reset link.",
	'forgot.sent':
		"If an account exists for that address, we've sent a link to reset your password. Check your inbox.",
	'forgot.submit': 'Send reset link',
	'forgot.backToLogin': 'Back to log in',

	// Reset password (set a new one).
	'reset.title': 'Set a new password - Bibseller',
	'reset.heading': 'Set a new password',
	'reset.done':
		"Your password has been updated. You've been signed out everywhere - sign in with your new password.",
	'reset.missingToken': 'This reset link is missing its token. Request a new one.',
	'reset.requestLink': 'Request a reset link',
	'reset.newPassword': 'New password',
	'reset.confirmPassword': 'Confirm password',
	'reset.submit': 'Update password',

	// Verify email (landing page).
	'verify.title': 'Verify email - Bibseller',
	'verify.okHeading': 'Email verified',
	'verify.okBody': "Your email address is confirmed - you're all set.",
	'verify.continue': 'Continue',
	'verify.invalidHeading': 'Link invalid or expired',
	'verify.invalidBody':
		'This verification link is no longer valid. Sign in and request a fresh one.',
	'verify.signIn': 'Sign in',
	'verify.missingHeading': 'Nothing to verify',
	'verify.missingBody': 'Open the verification link from your email to confirm your address.',
	'verify.home': 'Home',
	'verify.errorHeading': 'Something went wrong',
	'verify.errorBody': "We couldn't verify your email right now. Please try again in a moment.",

	// Settings.
	'settings.title': 'Settings - Bibseller',
	'settings.heading': 'Settings',
	'settings.navAria': 'Settings sections',
	'settings.security': 'Security',
	'settings.account': 'Account',
	'settings.profileHint': 'Your public name, language and country.',
	'settings.securityHint': 'Password and signed-in devices.',
	'settings.accountHint': 'Delete your account and data.',
	'settings.profile': 'Profile',
	'settings.country': 'Country',
	'settings.countryNotSet': 'Not set',
	'settings.profileUpdated': 'Profile updated.',
	'settings.save': 'Save',
	'settings.password': 'Password',
	'settings.currentPassword': 'Current password',
	'settings.confirmNewPassword': 'Confirm new password',
	'settings.passwordChanged': 'Password changed. Other devices have been signed out.',
	'settings.changePassword': 'Change password',
	'settings.sessions': 'Sessions',
	'settings.sessionsNote': 'Sign out of Bibseller on every device, including this one.',
	'settings.logoutAll': 'Log out all devices',
	'settings.deleteAccount': 'Delete account',
	'settings.deleteNote':
		'Permanently delete your account and its data. Available once full GDPR tooling ships (M7).',
	'settings.deleteTitle': 'Account deletion arrives with trust and safety (M7)',
	'settings.deleteSoon': 'Delete account - coming soon',

	// Error page.
	'error.notFound': "That page doesn't exist.",
	'error.generic': 'Something went wrong.',
	'error.backHome': 'Back home',

	// API errors (#49): keyed by the Go envelope's stable `code`, not its English
	// message, via apiErrorKey(). Any code without an entry falls back to
	// apiError.unknown - so no English API string is ever fed to the translator.
	'apiError.unknown': 'Something went wrong. Please try again.',
	'apiError.unreachable': 'The API is unreachable.',
	'apiError.not_found': 'Not found.',
	'apiError.invalid_parameter': 'Invalid request.',
	'apiError.email_taken': 'An account with this email already exists.',
	'apiError.race_past': 'This race has already taken place.',
	'apiError.loadFailed': 'Could not load this page. Please try again.',

	// Form/validation errors authored on the frontend (server actions). These are
	// our own copy, not the API's - the API re-validates and is the authority.
	'formError.invalidEmail': 'Enter a valid email address.',
	'formError.displayNameLength': 'Display name must be between 2 and 50 characters.',
	'formError.passwordTooShort': 'Password must be at least 8 characters.',
	'formError.newPasswordTooShort': 'New password must be at least 8 characters.',
	'formError.passwordMismatch': 'The two passwords do not match.',
	'formError.newPasswordMismatch': 'The two new passwords do not match.',
	'formError.loginRequired': 'Enter your email and password.',
	'formError.invalidCredentials': 'Invalid email or password.',
	'formError.loginFailed': 'Could not log in.',
	'formError.currentPasswordWrong': 'Your current password is incorrect.',
	'formError.changePasswordFailed': 'Could not change your password.',
	'formError.resetTokenMissing': 'This reset link is missing its token.',
	'formError.resetTokenInvalid': 'This reset link is invalid or has expired.',
	'formError.resetFailed': 'Could not reset your password.',
	'formError.resetEmailFailed': 'Could not send the reset email. Please try again.',
	'formError.tooManyRequests': 'Too many requests. Please wait a minute and try again.',
	'formError.missingRace': 'Missing race - please start again from the race page.',
	'formError.missingListingId': 'Missing listing id.',
	'formError.invalidAmount': 'Enter a valid amount, e.g. 45 or 45.00.',
	'formError.priceExceedsFace': 'Asking price cannot exceed the original face value.',
	'formError.verifyToContact': 'Verify your email before contacting a seller.',
	'formError.verifyToPublish': 'Verify your email before publishing a listing.',
	'formError.emptyMessage': 'Write a message to the seller first.',
	'formError.ackRequired': 'Please acknowledge the terms to continue.',
	'formError.ackFailed': 'Could not record your acknowledgment. Try again.',
	'formError.cannotContact': 'You cannot contact the seller for this listing.',
	'formError.listingUnavailable': 'This listing is no longer available.',
	'formError.contactFailed': 'Could not start the conversation.',
	'formError.editOwnOnly': 'You can only edit your own listing.',
	'formError.editNotActive': 'This listing is no longer active and cannot be edited.',
	'formError.editFailed': 'Could not update the listing.',
	'formError.cancelFailed': 'Could not cancel the listing.',

	// Legal pages (M7-lite, #11). DRAFT copy pending counsel review (D10); the
	// (legal) route group renders these under a draft banner.
	'legal.draftNotice': 'Draft - pending legal review and not yet binding.',
	'legal.terms.title': 'Terms of Service - Bibseller',
	'legal.terms.metaDescription': 'The terms governing use of the Bibseller marketplace.',
	'legal.terms.heading': 'Terms of Service',
	'legal.terms.intro':
		'Bibseller is a zero-commission marketplace that connects runners who want to pass on a race bib with runners looking for one. We provide the channel; we are not a party to any transfer between users.',
	'legal.terms.roleHeading': 'Our role',
	'legal.terms.roleBody':
		'We never charge commission and may only pass through payment-processor fees where secure payment is offered. What you may do with a bib is set by each race organiser, and the platform applies those rules per race.',
	'legal.terms.policyHeading': 'Responsibilities by race transfer policy',
	'legal.terms.policyIntro':
		'Every race carries a transfer policy that determines what the platform does and who is responsible for what:',
	'legal.terms.policy.platform_sale':
		"Resale allowed: the race permits transfers. You arrange the details in chat and, once secure payment ships, may pay through the platform; the bib transfer and race entry stay subject to the organiser's process.",
	'legal.terms.policy.official_only':
		'Official process: the race runs its own change-of-name procedure. We connect you and link to it; the transfer and any official fees happen with the organiser, never through the platform.',
	'legal.terms.policy.connect_only':
		"Connect only: the race restricts transfers. The platform only introduces you - it handles no money and takes no part in any agreement. The race's rules apply, and you confirm you understand this before contacting a seller.",
	'legal.terms.policy.unknown':
		"Unverified: until we confirm a race's policy we treat it as connect-only - introduction only, no money through the platform.",
	'legal.terms.userHeading': 'Your responsibilities',
	'legal.terms.userBody':
		"List bibs accurately, follow each race's transfer rules, and never ask above a bib's original face value. Do not use Bibseller for anything unlawful.",
	'legal.terms.conductHeading': 'Prohibited conduct and repeat infringers',
	'legal.terms.conductBody':
		"We may remove listings or messages and restrict or close accounts that break these terms or a race's rules. Accounts that repeatedly infringe will be terminated.",
	'legal.terms.liabilityHeading': 'Liability',
	'legal.terms.liabilityBody':
		'Except where the platform processes a payment for a resale-allowed race, transactions are strictly between users and the platform is not responsible for them. Nothing here excludes liability that cannot be excluded under applicable law.',
	'legal.terms.changesHeading': 'Changes and contact',
	'legal.terms.changesBody':
		'We may update these terms; material changes will be posted on this page. Questions go to our contact page.',

	// Privacy Policy (M7.2, #54). DRAFT pending counsel review (D10).
	'legal.privacy.title': 'Privacy Policy - Bibseller',
	'legal.privacy.metaDescription': 'How Bibseller handles your personal data.',
	'legal.privacy.heading': 'Privacy Policy',
	'legal.privacy.intro':
		'This policy explains what personal data Bibseller holds, why, and the rights you have over it.',
	'legal.privacy.dataHeading': 'What we hold',
	'legal.privacy.dataBody':
		'Your account (email, display name, preferred language, country), the bibs you list, the messages and images you exchange in chat, and any reports or blocks you make.',
	'legal.privacy.basisHeading': 'Why we hold it',
	'legal.privacy.basisBody':
		'To run the marketplace you asked to use: showing your listings, delivering your messages, and keeping the platform safe. We never sell your data.',
	'legal.privacy.retentionHeading': 'How long we keep it',
	'legal.privacy.retentionBody':
		'Account data stays while your account is open. Chat messages are removed on a rolling retention schedule. Some records may be kept longer where the law requires.',
	'legal.privacy.rightsHeading': 'Your rights',
	'legal.privacy.rightsBody':
		'You can access, correct, export, or delete your personal data. Self-serve export and deletion arrive with our full GDPR tooling; until then, contact us and we will action your request.',
	'legal.privacy.cookiesHeading': 'Cookies',
	'legal.privacy.cookiesBody':
		'We set one strictly necessary cookie to keep you signed in. We use no analytics or tracking cookies, so there is no consent banner.',
	'legal.privacy.contactHeading': 'Contact',
	'legal.privacy.contactBody':
		'For any privacy question or request, reach us through our contact page.',

	// Contact / imprint (M7.2, #54). Doubles as the DSA single point of contact.
	'legal.contact.title': 'Contact - Bibseller',
	'legal.contact.metaDescription':
		'How to contact Bibseller, including our point of contact for authorities.',
	'legal.contact.heading': 'Contact',
	'legal.contact.intro': 'How to reach Bibseller.',
	'legal.contact.operatorHeading': 'Who runs Bibseller',
	'legal.contact.operatorBody':
		'Bibseller is run as a non-profit marketplace. The operating entity, registered address, and a contact email will be published here before the public beta.',
	'legal.contact.dsaHeading': 'Point of contact (Digital Services Act)',
	'legal.contact.dsaBody':
		'This will also be our single point of contact for users and authorities under the EU Digital Services Act. You may write to us in English or Spanish.',
	'legal.contact.reportHeading': 'Reporting a listing',
	'legal.contact.reportBody':
		'To report an unlawful or rule-breaking listing, use the Report button on that listing - we review every report.'
} as const;

export type MessageKey = keyof typeof en;
