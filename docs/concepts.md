# Wiz Concepts

Here we define some terms as they are used in the context of the Wiz platform.

## Package
A package is a
- runtime
- framework/toolkit/library
- driver
- or ML model

A package defines it's type in it's config file.

## Component
A component is a piece of the final ML system that will be deployed. This can be a program or a package

## Dataset
A dataset is a collection of data

## Specification
A specification is a format/interface for
- An ML model
- or a dataset

Specifications can be defined in JSON or Protobuf

## ML Program
Well, what does it sound like? In this context, a program mainly
- trains
- runs
- or creates an ML model

A program can be written in any language, and can use any framework/library. The Wiz project provides a system for optimizing, extending, and distributing ML programs.
