
type MailHeaderKey = "From" | "To" | "Subject" | "Date" | "Content-Type" | "Return-Path";

type MailHeaders = {
    [key: MailHeaderKey | string]: string[];
};

export type MailData = {
    "id": string;
    "from": {
        "name": string;
        "domain": string;
    },
    "to": [
        {
            "name": string;
            "domain": string;
        }
    ],
    "created": string;
    "content": {
        "headers": MailHeaders;
        "size": number;
        "body": string;
    }
};

type MailSummary = {
  "emails": MailData[];
};

export default MailSummary;
