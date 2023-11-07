import React from "react";
import "../styles/post-categories.scss";
export const PostCategories = ({ categories }) => {
  return (
    <div className="post-categories">
      {categories?.map((category) => (
        <Category category={category} />
      ))}
    </div>
  );
};

const Category = ({ category }) => {
  return <div className="category">{category.name}</div>;
};
