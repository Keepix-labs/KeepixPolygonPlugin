import "./Faq.scss";

type Question = {
  title: string;
  desc: string;
};

const questions: Question[] = [
  {
    title: "How much does it cost ?",
    desc: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Tempora, qui expedita totam voluptatem consectetur dolores vel neque pariatur. Consequuntur veniam repellendus minus saepe culpa maiores, reprehenderit tempora perferendis quo? A!",
  },
  {
    title: "How much does it cost ?",
    desc: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Tempora, qui expedita totam voluptatem consectetur dolores vel neque pariatur. Consequuntur veniam repellendus minus saepe culpa maiores, reprehenderit tempora perferendis quo? A!",
  },
  {
    title: "How much does it cost ?",
    desc: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Tempora, qui expedita totam voluptatem consectetur dolores vel neque pariatur. Consequuntur veniam repellendus minus saepe culpa maiores, reprehenderit tempora perferendis quo? A!",
  },
  {
    title: "How much does it cost ?",
    desc: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Tempora, qui expedita totam voluptatem consectetur dolores vel neque pariatur. Consequuntur veniam repellendus minus saepe culpa maiores, reprehenderit tempora perferendis quo? A!",
  },
];

export default function FAQ({ questions }: any) {
  return (
    <div className="card card-default">
      <h2 className="h2">FAQ</h2>
      <dl className="questions">
        {questions.map((question: any) => (
          <details className="question" key={question.title}>
            <summary className="questionTitle">
              <span>{question.title}</span>
            </summary>
            <div className="questionDesc">{question.desc}</div>
          </details>
        ))}
      </dl>
    </div>
  );
}