`include "full_adder.sv"

module alu_adder_half (
   input  logic [3:0] a, b,
   input  logic       carry_in,
   input  logic       enable,
   output logic [3:0] result,
   output logic       carry_out,
);

    logic [3:0] internal_result;
    logic       internal_carry;

    full_adder #(.N(4)) adder(
        .a(a), .b(b),
        .carry_in(carry_in),
        .result(internal_result),
        .carry_out(internal_carry)
    );

    always_comb begin


        if (enable)
            result <= internal_result;
        else
            result <= 0;

        carry_out <= internal_carry;
    end

endmodule